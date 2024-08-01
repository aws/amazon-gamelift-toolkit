
# GameLift Plugin / C++ Server SDK Helper Script

This Python script automates the setup process for the Amazon GameLift Unreal Plugin. It handles the following tasks:

Downloads the latest version of the Amazon GameLift Plugin for Unreal Engine

- Verifies that the system meets the necessary installation requirements, including the presence of OpenSSL and other dependencies
- Builds the GameLift Unreal Server SDK and the GameLift C++ Server SDK
- Sets up resources, OpenSSL binaries, and dependencies in the correct project structure
- Handles pre-compiled library binaries and dependencies for specific platforms (if available)
- Performs final checks to ensure all components are correctly set up and ready for server deployment

Before running the script, ensure your system meets the following minimum requirements:
- ## Minimum requirements:
- ## Please note: If Strawberry Perl for Windows is installed, this might have conflics with these requirements.

* Either of the following:
    * Microsoft Visual Studio 2012 or later
    * GNU Compiler Collection (GCC) 4.9 or later
* CMake version 3.1 or later
* A Git client available on the PATH.
* OpenSSL installation, with the version matching the one used by Unreal Engine

### OpenSSL Installation
On Windows and Linux you can check to see if or which OpenSSL version is installed CLI $: openssl version

**Windows**:

1. Install the full version of OpenSSL from [Binaries - OpenSSLWiki](https://wiki.openssl.org/index.php/Binaries)
2. Build from source for Unreal OpenSSL 1.1.1  https://openssl-library.org/source/old/1.1.1/index.html
3. Unreal OpenSSL Unreal 5.2 - 5.1 (OpenSSL 1.1.1n): https://openssl.org/source/old/1.1.1/openssl-1.1.1n.tar.gz
4. Unreal OpenSSL Unreal 5.4 - 5.3 (OpenSSL 1.1.1t): https://github.com/openssl/openssl/releases/download/OpenSSL_1_1_1t/openssl-1.1.1t.tar.gz
5. Unreal OpenSSL Unreal 5.0 (OpenSSL 1.1.1c): https://openssl.org/source/old/1.1.1/openssl-1.1.1c.tar.gz
7. Add the OpenSSL install directory to the system PATH
8. Create the following environment variables:
   - `OPENSSL_INCLUDE_DIR = <PATH_TO_OPENSSL_DIR>\include`
   - `OPENSSL_LIBRARIES = <PATH_TO_OPENSSL_DIR>\lib`
   - `OPENSSL_ROOT_DIR = <PATH_TO_OPENSSL_DIR>\OpenSSL`

**Linux**: Example of download and building, use the version that matches your Unreal Version.

1. Unreal OpenSSL Unreal 5.2 - 5.1 (OpenSSL 1.1.1n): wget https://www.openssl.org/source/openssl-1.1.1n.tar.gz
3. Unreal OpenSSL Unreal 5.4 - 5.3 (OpenSSL 1.1.1t): wget https://www.openssl.org/source/openssl-1.1.1t.tar.gz
4. Unreal OpenSSL Unreal 5.0 (OpenSSL 1.1.1c): wget https://www.openssl.org/source/openssl-1.1.1c.tar.gz
5. Extract the source code (Example): tar -xzf openssl-1.1.1t.tar.gz
6. Change Directory (Example): cd openssl-1.1.1t
7. Configure and build OpenSSL 1.1.1 (Example): ./config --prefix=/usr/local/openssl-1.1.1
8. Make: make
9. Install: sudo make install
10. Set the environment variables (Example):
   export OPENSSL_INCLUDE_DIR=/usr/local/openssl-1.1.1/include
   export OPENSSL_LIBRARIES=/usr/local/openssl-1.1.1/lib
   export OPENSSL_ROOT_DIR=/usr/local/openssl-1.1.1
Then, reload the shell configuration: source ~/.bashrc
Verify the installation: openssl version

## Usage

1. Download the `GLDownloadPluginBuildSDK.py` script.

2. Run the script with the appropriate command for your platform:
   
   **Windows and Linux**:
   
   ```
   python GLDownloadPluginBuildSDK.py
   ```
   
   or
   
   ```
   python3 GLDownloadPluginBuildSDK.py
   ```

## Support

Please report issues under issues in this repository


3. Open a terminal or command prompt and navigate to the script's directory.

## Output:
This will be the full plugin with the built SDKs

Windows:
 C:\Users\{USERNAME}\.cache\AmazonGameLift\amazon-gamelift-plugin-unreal\GameLiftPlugin

Linux:
$home\{USERNAME}\.cache\AmazonGameLift\amazon-gamelift-plugin-unreal\GameLiftPlugin
