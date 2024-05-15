# Golang WebAuthn Example

This is a fully functional Golang example for logging in with WebAuthn.

## How to Run

1. Clone or Download this repository;
2. Type `go run .` to run the server;
3. Go to <https://localhost:2345>.

You should make your browser trust this self-signed certificates.

## How to Play

1. Go to `/admin/`;
2. Sign up with any valid email address;
3. Add you WebAuthn Credential;
4. Upon successful add, log out;
5. Try to use your newly created credential to log in again.

## To generate new Certificate?

Just run `keygen.sh`. Only the `FQDN` field is required.

## File structure

TBD.

## References

* [Guide to Web Authentication](https://webauthn.guide/)
* [WebAuthn.io](https://webauthn.io/)
* [Passkeys: the web authentication standard](https://www.passkeys.com/)
* [Meet passkeys - WWDC22 - Videos - Apple Developer](https://developer.apple.com/videos/play/wwdc2022/10092/)
* [Introduction to WebAuthn API and Passkey](https://medium.com/webauthnworks/introduction-to-webauthn-api-5fd1fb46c285)
* [Apple Just Killed the Password—for Real This Time](https://www.wired.com/story/apple-passkeys-password-ios16-ventura/)
* [MacOS 13 Ventura: Features, Details, Release Date | WIRED](https://www.wired.com/story/apple-ventura-macos-13-preview/)
* [Web Authentication API - Web APIs | MDN](https://developer.mozilla.org/en-US/docs/Web/API/Web_Authentication_API)
* [Passkeys in iOS 16 and macOS 13 enable passwordless sign-in - 9to5Mac](https://9to5mac.com/2022/06/07/passkeys-passwordless-sign-in-ios-16/)
* [duo-labs/webauthn: WebAuthn (FIDO2) server library written in Go](https://github.com/duo-labs/webauthn)
* [duo-labs/webauthn.io: The source code for webauthn.io, a demonstration of WebAuthn.](https://github.com/duo-labs/webauthn.io)
* [koesie10/webauthn: Go package for easy WebAuthn integration](https://github.com/koesie10/webauthn)
* [hbolimovsky/webauthn-example: Basic WebAuthn client and server in go](https://github.com/hbolimovsky/webauthn-example)
* [WebAuthn Basic Web Client/Server · Herbie's Blog](https://www.herbie.dev/blog/webauthn-basic-web-client-server/)
