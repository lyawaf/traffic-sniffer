#!/usr/bin/env python3


from pyshark import FileCapture as pcap
from modules.message import tcpMessage
from modules.config import SharkTcpConfig
from modules.constants import SharkTcpConstants


class SharkTcpStreamer:
    """
    Pcap TCP stream id parser
    """

    def __init__(self, filename):
        """
        Initialize class object with current
        filename of pcap file

        :param filename: name of pcap file for parsing
        """
        self.filename: str = filename
        self.capture = pcap(filename)
        self.stream_dict: dict = {}

    def decode_from_hex(self, payload):
        """
        Remove all special characters from payload

        :param payload: wireshark packet payload
        :return payload: decoded packet payload
        """

        # Ignore any non-printable characters errors
        payload = bytearray.fromhex(payload).decode("utf-8", errors="ignore")

        # Remove special characters like \x01 etc.
        if SharkTcpConfig.SET_REMOVE_CONTROL_CHARACTERS:
            payload = payload.translate(dict.fromkeys(range(32)))

        return payload

    def clear_hex(self, payload):
        """
        Remove ":" from packet payload

        :param payload: wireshark packet payload
        :return payload: payload without splitters
        """
        return payload.replace(":", "")

    def parse_by_stream(self):
        """
        Parse pcap file in dict of streams
        (stream_id: stream_info)

        :return:
        """
        server_port = None
        stream_time = None

        for packet in self.capture:
            # Reject all packets without TCP layers
            if not "TCP" in packet:
                continue

            # If we got SYN flag in tcp packet,
            # our server port = destination port (connect state)
            if packet.tcp.flags == "0x00000002":
                server_port = packet[packet.transport_layer].dstport

            # Debug output for current stream id
            stream = int(packet.tcp.stream)
            if SharkTcpConfig.SET_DEBUG_MODE:
                print(f"Current stream: {stream}")

            # Check if we got payload in tcp packet
            # or it is empty
            tcp_payload = packet.tcp.get("payload")
            if tcp_payload:
                payload = self.clear_hex(tcp_payload)

            # Detect direction of messages
            if packet[packet.transport_layer].dstport == server_port:
                direction = SharkTcpConstants.DIR_CLIENT_SERVER
            elif packet[packet.transport_layer].srcport == server_port:
                direction = SharkTcpConstants.DIR_SERVER_CLIENT
            else:
                direction = SharkTcpConstants.DIR_UNKNOWN

            # Build our message
            current_message = tcpMessage(
                direction=direction,
                data=payload if tcp_payload else None,
                time=packet.sniff_time,
            )

            # And put it in list for easy extending/appending
            payload_list = [current_message]

            # Build up our dictionary
            if not self.stream_dict.get(stream):
                self.stream_dict.update({stream: {}})

            # Set or extend current message on stream
            if not self.stream_dict[stream].get("messages"):
                self.stream_dict[stream].update({"messages": payload_list})
            else:
                self.stream_dict[stream]["messages"].extend(payload_list)

            stream_packet_values = {
                "filename": self.filename,
                "port": server_port,
                "direction": direction,
                "time": stream_time,
            }

            for packet_item in stream_packet_values.items():
                if not self.stream_dict[stream].get(packet_item[0]):
                    self.stream_dict[stream].update({packet_item[0]: packet_item[1]})

    def get_list_of_streams(self):
        """
        Get list of all streams with unique stream info

        :return: list of streams
        """
        self.stream_list = []
        for stream_id in self.stream_dict.keys():
            stream_info = {
                "filename": self.stream_dict[stream_id]["filename"],
                "stream_id": stream_id,
                "port": self.stream_dict[stream_id]["port"],
                "messages": self.stream_dict[stream_id]["messages"],
            }
            self.stream_list.append(stream_info)
        return self.stream_list
