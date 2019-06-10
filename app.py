from flask import Flask, render_template, redirect, url_for, request


class DefaultCredentials:
    username = "admin"
    password = "password"


app = Flask(__name__)


@app.route("/")
def homepage():
    return redirect(url_for("welcome"))


@app.route("/welcome")
def welcome():
    return render_template("welcome.html")


@app.route("/success")
def success_login():
    return render_template("successful.html", username=DefaultCredentials.username)


@app.route("/login", methods=["GET", "POST"])
def login():
    if request.method == "POST":
        if (request.form.get("username"), request.form.get("password")) == (
            DefaultCredentials.username,
            DefaultCredentials.password,
        ):
            return redirect(url_for("success_login"))
        return render_template(
            "login.html", error="Invalid Credentials. Please try again."
        )
    return render_template("login.html")


if __name__ == "__main__":
    app.run(debug=True)
