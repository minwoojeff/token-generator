import sys
import requests


def health_check(url):
    health_url = '{}/health'.format(url)

    r = requests.get(health_url)
    assert(r.status_code == 200)


def authentication_validation(url):
    oauth_token_url = '{}/oauth/token'.format(url)

    no_auth(oauth_token_url)
    basic_auth(oauth_token_url)
    body_auth(oauth_token_url)


def no_auth(url):
    r = requests.post(url, json={})
    assert(r.status_code == 401)


def basic_auth(url):
    r = requests.post(url, auth=("default", "password"), json={})
    assert(r.status_code == 400)
    assert(r.text == 'Unsupported grant type')


def body_auth(url):
    payload = {
        "client_id": "default",
        "client_secret": "password"
    }
    r = requests.post(url, json=payload)
    assert(r.status_code == 400)
    assert(r.text == 'Unsupported grant type')


def main():
    args = sys.argv

    if len(args) != 2:
        print("Invalid number of args. Aborting")
        exit(0)

    base_url = args[1]
    print("URL: {}".format(base_url))

    # 1. Health check
    health_check(base_url)
    print("PASS: Health Check")

    # 2. Authentication Validation
    authentication_validation(base_url)
    print("PASS: Authentication")


if __name__ == "__main__":
    main()
