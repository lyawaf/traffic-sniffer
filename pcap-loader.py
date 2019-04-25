#!/usr/bin/env python3

import os
from pprint import PrettyPrinter
from modules.config import SharkTcpConfig
from modules.constants import SharkTcpConstants
from modules.parser import SharkTcpStreamer

SharkTcpConfig.SET_DEBUG_MODE = True
prettyp = PrettyPrinter(indent=4)


def get_pcap_files():
    """
    Get list of all *.pcap files from `pcap` dir
    :return: list of pcap files with full paths
    """
    pcap_list = []
    for file in os.listdir("pcap"):
        if file.endswith(".pcap"):
            pcap_list.append(os.path.join("pcap", file))
    return pcap_list


def get_stream_list(pcap):
    """
    Get list of streams with info for current
    pcap file

    :param pcap: name of pcap file
    :return: list of streams in pcap file
    """
    streamer = SharkTcpStreamer(pcap)
    streamer.parse_by_stream()
    stream_list = streamer.get_list_of_streams()
    return stream_list


def main():
    """
    For testing purpose

    :return:
    """
    pcap_files = get_pcap_files()
    pcap_stream_list = []

    # Currently - one pcap
    for pcap in pcap_files:
        stream_list = get_stream_list(pcap)
        pcap_stream_list.append(stream_list)

    for stream_list in pcap_stream_list:
        prettyp.pprint(stream_list)


if __name__ == "__main__":
    main()
