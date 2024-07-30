- 

# GameLift Plugin / C++ Server SDK Helper Script

This Python script automates the setup process for the Amazon GameLift Unreal Plugin. It handles the following tasks:

Downloads the latest version of the Amazon GameLift Plugin for Unreal Engine

- Verifies that the system meets the necessary installation requirements, including the presence of OpenSSL and other dependencies
- Builds the GameLift Unreal Server SDK and the GameLift C++ Server SDK
- Sets up resources, OpenSSL binaries, and dependencies in the correct project structure
- Handles pre-compiled library binaries and dependencies for specific platforms (if available)
- Performs final checks to ensure all components are correctly set up and ready for server deployment

## Prerequisites

Before running the script, ensure your system meets the following minimum requirements:

- Microsoft Visual Studio 2012 or later, or GNU Compiler Collection (GCC) 4.9 or later
- CMake version 3.1 or later
- A Git client available on the PATH
- OpenSSL installation, with the version matching the one used by Unreal Engine

### OpenSSL Installation

**Windows**:

1. Install the full version of OpenSSL from [Binaries - OpenSSLWiki](https://wiki.openssl.org/index.php/Binaries) or use pre-built ones: [OpenSSL_1.1.1n-win64-SS.zip](https://wiki.agscollab.com/download/attachments/206439670/OpenSSL_1.1.1n-win64-SS.zip?version=1&modificationDate=1721847618073&api=v2)
2. Add the OpenSSL install directory to the system PATH
3. Create the following environment variables:
   - `OPENSSL_INCLUDE_DIR = <PATH_TO_OPENSSL_DIR>\include`
   - `OPENSSL_LIBRARIES = <PATH_TO_OPENSSL_DIR>\lib`
   - `OPENSSL_ROOT_DIR = <PATH_TO_OPENSSL_DIR>\OpenSSL`

**Linux**:

1. Install the full version of OpenSSL from [Binaries - OpenSSLWiki](https://wiki.openssl.org/index.php/Binaries)
2. Follow the platform-specific instructions for setting up the OpenSSL environment variables

## Usage

1. Download the `GLDownloadPluginBuildSDK.py` script.

2. Open a terminal or command prompt and navigate to the script's directory.

3. Run the script with the appropriate command for your platform:
   
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

## Downloads:

1. Download the `GLDownloadPluginBuildSDK.py` script here > [![](https://wiki.agscollab.com/s/ubtom5/8804/1yuue1v/17.19.4/_/download/resources/com.atlassian.confluence.plugins.confluence-view-file-macro:view-file-macro-resources/images/placeholder-medium-file.png)GLDownloadPluginBuildSDK.py](https://wiki.agscollab.com/download/attachments/206439670/GLDownloadPluginBuildSDK.py?version=9&modificationDate=1721855344979&api=v2)

2. Internal Testing Only Get past AWS VPN on Linux: [GLDownloadPluginBuildSDK.py](https://wiki.agscollab.com/download/attachments/206439670/GLDownloadPluginBuildSDK.py?version=9&modificationDate=1721855344979&api=v2)
   
   v0.5 24 Jul 2024 

3. Open a terminal or command prompt and navigate to the script's directory.

## Output:

This will be the full plugin in your C:\Users\{USERNAME}\.cache\AmazonGameLift\amazon-gamelift-plugin-unreal\GameLiftPlugin

![](https://wiki.agscollab.com/download/attachments/206439670/image-2024-7-24_16-28-21.png?version=1&modificationDate=1721856501190&api=v2)

Windows: 

![](https://wiki.agscollab.com/download/attachments/206439670/image-2024-7-24_16-28-52.png?version=1&modificationDate=1721856532231&api=v2)

Linux:

![](https://wiki.agscollab.com/download/attachments/206439670/image-2024-7-24_16-29-48.png?version=1&modificationDate=1721856588328&api=v2)
