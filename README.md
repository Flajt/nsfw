# NSFW

(Network Security File Writer)
<br>

Generates you Network Security files for Android (network_security_config.xml) and the NSAppTransportSecurity part iOS via a simple cli.

> [!NOTE]
> Info: Note that the Info.plist might not be "valid", since there is no header and it starts with a key tag

- 📦 No external dependencies
- 📁 Output as file or as plain text
- 📅 Fetch certificate expiry time automaticaly
- 📌 Pin multiple websites in one go

## Download Binary

- TODO

## Building
1. Clone the repo
2. Build using `go build`
3. Run

## Commands
- `-websites`: Takes in a comma seperated list of domains e.g. https://www.test.com,https://test2.com (**required**)
- `-output`: Takes in the output path for generated files (**optional**)
- `-no-file`: Writes data to `stdout` if set, otherwise to `Info.plst` & `networks_security_config.xml` (**optional**)
- `-platforms`: Specify for which platforms you want to generate the pinning files, comma sperated e.g. android,ios (**optional**)
- `-help`: Returns help interface

### Example

- `nsfw -websites https://test.com -output ~/my-app/android/app/src/main/xml/ressources -platforms android`: Creates the android `network_security_config.xml` in its folder for test.cmo
- `nsfw -websites https://test.com,https://www.test2.com -no-file -platforms ios`: Returns the `NSNetworkAppTransportSecurity` config in the terminal
- `nsfw -websites https://test.com https://www.test2.com`: Generate both config files for ios and android and saves them into the current folder

## Not supported

- Providing fallbacks
- Providing multiple pins
- Provide certificates
- Disabling subdomain pinning support (May be possible to implement if desired)
- Pin more than the leaf certificate (May be possible to implement if desired)
- Customize file names (May be possible to implement if desired)
- Custom output path for each android & ios (Possible if desired)

## TODOS:

- Figure out if I follow best practices, since I'm new to go
- Review code organization
- Consider moving networking into sperate struct / function
- Running requests in paralel?
- Finding a way to embed / merge the ios output into existing Info.plist file 
- Website, so you can run everything via a comfy ui
- Testing?
- Handle panics gracefully (not throwing around stack traces etc.)
- Setup Github action for releasing binarys
- (Packageing code for easier installation)
- (Offering CI/CD integrations)

## Contributions

Feel free to either pick something from the TODOS, or open your own issue and we can discuss things<br>

Feedback for the existing codebase is also highly appreciated!
