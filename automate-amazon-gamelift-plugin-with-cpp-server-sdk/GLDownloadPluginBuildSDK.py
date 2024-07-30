# CLI Command
# python3 GLDownloadPluginBuildSDK.py

import sys
import subprocess
import zipfile
import shutil
from pathlib import Path
import os
import logging
import webbrowser
import re
import platform
import urllib.request


is_windows = platform.system() == "Windows"
is_linux = platform.system() == "Linux"
is_mac = platform.system() == "Darwin"  # Darwin is the kernel for macOS

if is_windows:
    operating_system = "Win64"
    print("Running on Windows")
    import winreg
elif is_linux:
    operating_system = "Linux"
    print("Running on Linux")
elif is_mac:
    operating_system = "OSX"
    print("Running on macOS")
else:
    print("Running on an unknown operating system")

# Set up logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

if operating_system == "Win64":
    import winreg
    CURRENT_OS_LOGIN = os.getlogin()

# CONSTANTS
VERSION = "v0.5 5/20/24"
MIN_CMAKE_VERSION = "3.1"
FILE_CACHE_PATH = Path.home() / '.cache'
OSX_PLUGIN_BUILD_FOLDER = Path.home() / 'gl_plugin_build_osx'
if operating_system == "OSX":
    if not OSX_PLUGIN_BUILD_FOLDER.exists():
        OSX_PLUGIN_BUILD_FOLDER.mkdir(parents=True)
    FILE_CACHE_PATH = OSX_PLUGIN_BUILD_FOLDER
CURRENT_GAMELIFT_ZIP = "amazon-gamelift-plugin-unreal-release-1.1.1.zip"
CURRENT_GAMELIFT_CPP_SDK_NAME = "GameLift-Cpp-ServerSDK-5.1.2"
GAMELIFT_PLUGIN_DOWNLOAD_URL = f"https://github.com/aws/amazon-gamelift-plugin-unreal/releases/download/v1.1.1/{CURRENT_GAMELIFT_ZIP}"
GAMELIFT_PLUGIN_DOWNLOAD_PATH = Path(FILE_CACHE_PATH, CURRENT_GAMELIFT_ZIP)
GAMELIFT_PLUGIN_EXTRACT_PATH = Path(FILE_CACHE_PATH, "AmazonGameLift")
GAMELIFT_PLUGIN_MAIN_ZIP_PATH = Path(FILE_CACHE_PATH, "AmazonGameLift", "amazon-gamelift-plugin-unreal-1.1.1-sdk-5.1.1.zip")
GAMELIFT_CCP_SERVER_SDK_ZIP_PATH = Path(FILE_CACHE_PATH, "AmazonGameLift", "GameLift-Cpp-ServerSDK-5.1.2.zip")
GAMELIFT_PLUGIN_MAIN_PATH = Path(FILE_CACHE_PATH, "AmazonGameLift")
GAMELIFT_SERVER_SDK_MAIN_PATH = Path(FILE_CACHE_PATH, "AmazonGameLift", "GameLift-Cpp-ServerSDK-5.1.2")
GAMELIFT_SERVER_SDK_BUILD_PATH = Path(FILE_CACHE_PATH, "AmazonGameLift", CURRENT_GAMELIFT_CPP_SDK_NAME)
GAMELIFT_SERVER_SDK_CMAKE_BUILD_PATH = Path(FILE_CACHE_PATH, "AmazonGameLift", CURRENT_GAMELIFT_CPP_SDK_NAME, "cmake-build")
SDK_DLL_BUILD_PATH = Path(FILE_CACHE_PATH, "AmazonGameLift", CURRENT_GAMELIFT_CPP_SDK_NAME, "cmake-build", "prefix", "bin", "aws-cpp-sdk-gamelift-server.dll")
SDK_LIB_BUILD_PATH = Path(FILE_CACHE_PATH, "AmazonGameLift", CURRENT_GAMELIFT_CPP_SDK_NAME, "cmake-build", "prefix", "lib", "aws-cpp-sdk-gamelift-server.lib")
SDK_SO_BUILD_PATH = Path(FILE_CACHE_PATH, "AmazonGameLift", CURRENT_GAMELIFT_CPP_SDK_NAME, "cmake-build", "prefix", "lib", "libaws-cpp-sdk-gamelift-server.so")
GAMELIFT_SERVER_SDK_PATH_WIN = Path(FILE_CACHE_PATH, "AmazonGameLift", "amazon-gamelift-plugin-unreal", "GameLiftPlugin", "Source", "GameLiftServer", "ThirdParty", "GameLiftServerSDK", "Win64")
GAMELIFT_SERVER_SDK_PATH_LINUX = Path(FILE_CACHE_PATH, "AmazonGameLift", "amazon-gamelift-plugin-unreal", "GameLiftPlugin", "Source", "GameLiftServer", "ThirdParty", "GameLiftServerSDK", "Linux", "x86_64-unknown-linux-gnu")
GAMELIFT_CCP_SERVER_SDK_ZIP_URL = "https://gamelift-server-sdk-release.s3.us-west-2.amazonaws.com/cpp/GameLift-Cpp-ServerSDK-5.1.2.zip"
CMAKE_URL = "https://cmake.org/download/"
GIT_URL = "https://git-scm.com/downloads/"
OPEN_SSL_URL = "https://wiki.openssl.org/index.php/Binaries/"
GL_DOCS_URL = "https://aws.amazon.com/documentation/gamelift/"

# Create the cache directory if it doesn't exist
if not FILE_CACHE_PATH.exists():
    try:
        FILE_CACHE_PATH.mkdir(parents=True)
    except PermissionError:
        # Handle the case where the user doesn't have permission to create the directory
        logging.error(f"Unable to create directory {FILE_CACHE_PATH}. Please check permissions.")
        sys.exit(1)

VISUAL_STUDIO = {
    "2022": "Visual Studio 17 2022",
    "2019": "Visual Studio 16 2019",
    "2017": "Visual Studio 15 2017",
    "2015": "Visual Studio 14 2015",
    "2013": "Visual Studio 12 2013",
    "2012": "Visual Studio 11 2012"
}


'''
def download_file(url: str, file_path: Path):
    """
    Downloads a file from the given URL and saves it to the specified file path.

    Args:
        url (str): URL download string
        file_path (Path): download .zip from AWS
    """
    response = ""
    try:
        # Lets use requests to download the .zip file
        with urllib.request.urlopen(url) as response:
            with open(file_path, 'wb') as file:
                file.write(response.read())
        logging.info(f"(DOWNLOAD) File downloaded successfully from {url} and saved to {file_path}")
    except urllib.error.URLError as e:
        logging.error("(DOWNLOAD) Issue downloading the GameLift Plugin .zip release")
        logging.error(f"Error downloading file: {e.reason}")
        logging.error("(DOWNLOAD) Do you have a work VPN on?")
        sys.exit(1)
'''
import ssl

def download_file(url: str, file_path: Path):
    """
    Downloads a file from the given URL and saves it to the specified file path.

    Args:
        url (str): URL download string
        file_path (Path): download .zip from AWS
    """
    try:
        # Create an unverified context
        context = ssl._create_unverified_context()

        # Use the unverified context when opening the URL
        with urllib.request.urlopen(url, context=context) as response:
            with open(file_path, 'wb') as file:
                file.write(response.read())
        logging.info(f"(DOWNLOAD) File downloaded successfully from {url} and saved to {file_path}")
    except urllib.error.URLError as e:
        logging.error("(DOWNLOAD) Issue downloading the GameLift Plugin .zip release")
        logging.error(f"Error downloading file: {e.reason}")
        logging.error("(DOWNLOAD) Do you have a work VPN on?")
        sys.exit(1)

def set_open_ssl_env_osx() -> bool:
    """
    This function sets the OpenSSL environment variables for macOS

    Returns:
        bool: Return if we set the OpenSSL environment variables
    """
    try:
        # Run the shell command to get the path to the OpenSSL binary
        openssl_path = subprocess.check_output(['which', 'openssl']).decode().strip()
        # Get the directory containing the OpenSSL binary
        openssl_root_dir = subprocess.check_output(['dirname', openssl_path]).decode().strip()
        logging.info(f"OpenSSL root directory: {openssl_root_dir}")

        openssl_libraries = os.path.join(openssl_root_dir, "lib")

        # Set the environment variables
        os.environ["OPENSSL_ROOT_DIR"] = openssl_root_dir
        os.environ["OPENSSL_LIBRARIES"] = openssl_libraries

        if "OPENSSL_ROOT_DIR" in os.environ or "OPENSSL_LIBRARIES" in os.environ:
            logging.info("(OPENSSL) OpenSSL environment variables are set for macOS.")
            return True
        else:
            logging.error("(OPENSSL) OpenSSL environment variables are not set for macOS.")
            return False
    except:
        logging.error("Issue with OpenSSL env")
        return False

def check_openssl_env(operating_system: str) -> bool:
    """
    This function checks to see if OpenSSL environment variables are set

    Returns:
        bool: Returns true or false depending on if OpenSSL environment variables are set
    """
    logging.info("(OPENSSL) Checking OpenSSL environment variables...")
    # Check if environment variables are set for Windows
    if operating_system == "Win64":
        if "OPENSSL_INCLUDE_DIR" in os.environ and "OPENSSL_LIBRARIES" in os.environ and "OPENSSL_ROOT_DIR" in os.environ:
            logging.info("(OPENSSL) OpenSSL environment variables are set for Windows.")
            return True
        else:
            logging.error("(OPENSSL) OpenSSL environment variables are not set for Windows.")
            return False
    # Check if environment variables are set for macOS
    elif operating_system == "OSX":
        if "OPENSSL_ROOT_DIR" in os.environ or "OPENSSL_LIBRARIES" in os.environ:
            logging.info("(OPENSSL) OpenSSL environment variables are set for macOS.")
            return True
        else:
            env_set = set_open_ssl_env_osx()

            if env_set:
                return True
            else:
                logging.error("(OPENSSL) OpenSSL environment variables are not set for macOS.")
                return False
    # For Linux, no environment variables need to be set
    else:
        logging.info("(OPENSSL) No OpenSSL environment variables are required for Linux.")
        return True

def check_openssl(operating_system: str) -> bool:
    """
    This function checks to see if OpenSSL is installed

    Returns:
        bool: Returns if OpenSSL is installed or not
    """
    try:
        output = subprocess.run(['openssl', 'version'], stdout=subprocess.PIPE, check=True)
        if output.stdout:
            logging.info(f"(OPENSSL) OpenSSL version: {output.stdout.decode().strip()}")
            environment_var = check_openssl_env(operating_system)
            return environment_var
        else:
            logging.error("(OPENSSL) OpenSSL is not installed.")
            return False
    except Exception as e:
        logging.error(f"(OPENSSL) OpenSSL is not installed: {e}")
        return False

def check_for_visual_studio() -> str:
    """
    This function checks to see if Visual Studio is installed

    Returns:
        tuple: (visual_studio_installed_version, vs_studio_path, edition)
    """
    print("(VISUAL STUDIO) Checking Visual Studio is installed..")
    visual_studio_installed_version = ""

    try:
        editions = ["Community", "Enterprise"]

        for year in range(2022, 2010, -1):
            for ed in editions:
                install_path = Path(r"C:\Program Files\Microsoft Visual Studio", str(year), ed)
                if install_path.exists():
                    visual_studio_installed_version = str(year)
                    break
            else:
                continue
            break

        if not visual_studio_installed_version:
            visual_studio_installed_version = None

        return visual_studio_installed_version
    except (FileNotFoundError, PermissionError, OSError):
        print("(VISUAL STUDIO) Issue looking for Visual Studio is installed..")
        return visual_studio_installed_version

def check_cmake_version() -> bool:
    """ This function will check the cmake version to see if it meets the min required

    Returns:
        bool: True if cmake version meets min required, False otherwise
    """
    try:
        # Specify the full path to the cmake executable
        cmake_version = subprocess.check_output(["cmake", "--version"])

        if cmake_version:
            # Regular expression pattern to extract version number
            pattern = r'cmake version (\d+\.\d+\.\d+)'
            match = re.search(pattern, cmake_version.decode('utf-8'))

            # If a match is found, extract the version number
            if match:
                cmake_version = match.group(1)
                if cmake_version >= MIN_CMAKE_VERSION:
                    logging.info(f"(CMAKE) CMake version {cmake_version} is installed, which meets the minimum requirement.")
                return True
            else:
                logging.error(f"(CMAKE) Installed CMake version {cmake_version} does not meet the minimum requirement of {MIN_CMAKE_VERSION}.")
                return False

    except Exception as e:
        logging.error(f"(CMAKE) {e}")
        logging.info("(CMAKE) CMake must be installed and available on PATH")
        logging.info("(CMAKE) If you just installed cmake you may need to restart system.")
        return False

def check_git() -> bool:
    """
    This function checks to see if Git is installed

    Returns:
        bool: Returns if Git is installed or not
    """
    logging.info("(GIT) Checking Git is installed..")
    git_version = ""
    try:
        git_version = subprocess.check_output(["git", "--version"]).decode().split()[-1]
        logging.info(f"(GIT) Git version: {git_version}")
        return True
    except subprocess.CalledProcessError:
        logging.error("(GIT) Git must be installed and available on PATH")
        return False
    

def build_sdk(visual_studio_installed_version: str, operating_system: str):
    """
    This will build the GameLift Server SDK and output the binaries needed
    for the GameLift Plugin Build Successful: aws-cpp-sdk-gamelift-server.dll

    Args:
        visual_studio_installed_version (str): Installed version of Visual Studio on Windows
    """
    # Let's make a new directory called cmake-build using pathlib
    if not GAMELIFT_SERVER_SDK_CMAKE_BUILD_PATH.exists():
        Path(GAMELIFT_SERVER_SDK_CMAKE_BUILD_PATH).mkdir(parents=True, exist_ok=True)

    logging.info("(BUILDING SDK) Building SDK...")
    # Let's change the current working directory to the cmake-build directory
    if Path(GAMELIFT_SERVER_SDK_BUILD_PATH).exists():
        os.chdir(GAMELIFT_SERVER_SDK_BUILD_PATH)
    # Let's use subprocess to run cmake
    if operating_system == "Win64":
        cmd = ["cmake.exe", "-G", f"{visual_studio_installed_version}", "-A", "x64", "-DBUILD_FOR_UNREAL=1", "-DCMAKE_BUILD_TYPE=Release", "-S", ".", "-B", str(GAMELIFT_SERVER_SDK_CMAKE_BUILD_PATH)]
    elif operating_system == "Linux":
        cmd = ["cmake", "-G", "Unix Makefiles", "-DBUILD_FOR_UNREAL=1", "-DCMAKE_BUILD_TYPE=Release", "-S", ".", "-B", str(GAMELIFT_SERVER_SDK_CMAKE_BUILD_PATH)]
    elif operating_system == "OSX":
        cmd = ["cmake", "-G", "Unix Makefiles", "-DBUILD_FOR_UNREAL=1", "-DCMAKE_BUILD_TYPE=Release", "-S", ".", "-B", str(GAMELIFT_SERVER_SDK_CMAKE_BUILD_PATH)]

    try:
        output = subprocess.run(cmd, check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        logging.info(output.stdout.decode())
        logging.info("(BUILDING SDK) Building please wait......")
    except subprocess.CalledProcessError as e:
        logging.error("(BUILDING SDK) Issue building the GameLift MakeFiles")
        logging.error(e.stderr.decode())  # Decode the stderr output
        sys.exit(1)

    try:
        # Let's use subprocess to run cmake --build ./cmake-build --target ALL_BUILD
        if operating_system == "Win64":
            cmd = ["cmake.exe", "--build", str(GAMELIFT_SERVER_SDK_CMAKE_BUILD_PATH), "--target", "ALL_BUILD"]
        else:
            cmd = ["cmake", "--build", str(GAMELIFT_SERVER_SDK_CMAKE_BUILD_PATH), "--target", "all"]

        output = subprocess.run(cmd, check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        logging.info(output.stdout.decode())
    except subprocess.CalledProcessError as e:
        logging.error("(BUILDING SDK) Issue building the GameLift Server SDK")
        logging.error(e.stderr.decode())  # Decode the stderr output
        sys.exit(1)

    # Let's copy the lib files from the source to the destination
    if operating_system == "Win64":
        if GAMELIFT_SERVER_SDK_PATH_WIN.exists():
            if SDK_DLL_BUILD_PATH.exists():
                shutil.copy(SDK_DLL_BUILD_PATH, GAMELIFT_SERVER_SDK_PATH_WIN)
            else:
                logging.info("(BUILDING SDK) Issue building the GameLift Server SDK")
            if SDK_LIB_BUILD_PATH.exists():
                shutil.copy(SDK_LIB_BUILD_PATH, GAMELIFT_SERVER_SDK_PATH_WIN)
            else:
                logging.info("(BUILDING SDK) Issue building the GameLift Server SDK")

            if Path(GAMELIFT_SERVER_SDK_PATH_WIN, SDK_DLL_BUILD_PATH).exists():
                logging.info(f"(BUILDING SDK) Build Successful: aws-cpp-sdk-gamelift-server.dll")
                logging.info("Added to the GameLift Plugin")
            if Path(GAMELIFT_SERVER_SDK_PATH_WIN, SDK_LIB_BUILD_PATH).exists():
                logging.info(f"(BUILDING SDK) Build Successful: aws-cpp-sdk-gamelift-server.lib")
                logging.info("Added to the GameLift Plugin")   
        else:
            logging.error(f"(BUILDING SDK) Issue building the GameLift Server SDK do to file missing: {GAMELIFT_SERVER_SDK_PATH_WIN}")
            sys.exit(1)

    elif operating_system == "Linux":
        if Path(GAMELIFT_SERVER_SDK_PATH_LINUX).exists():
            shutil.copy(SDK_SO_BUILD_PATH, GAMELIFT_SERVER_SDK_PATH_LINUX)
        else:
            logging.info("(BUILDING SDK) Issue building the GameLift Server SDK")

        if Path(SDK_SO_BUILD_PATH).exists():
            logging.info(f"(BUILDING SDK) Build Successful: libaws-cpp-sdk-gamelift-server.so")
            logging.info("(BUILDING SDK) Added to the GameLift Plugin")
        else:
            logging.error(f"(BUILDING SDK) Issue building the GameLift Server SDK do to file missing: {GAMELIFT_SERVER_SDK_PATH_LINUX}")
            sys.exit(1)
        
    elif operating_system == "OSX":
        if Path(GAMELIFT_SERVER_SDK_PATH_LINUX).exists():
            shutil.copy(SDK_SO_BUILD_PATH, GAMELIFT_SERVER_SDK_PATH_LINUX)
        else:
            logging.info("(BUILDING SDK) Issue building the GameLift Server SDK")

        if Path(SDK_SO_BUILD_PATH).exists():
            logging.info(f"(BUILDING SDK) Build Successful: libaws-cpp-sdk-gamelift-server.so")
            logging.info("(BUILDING SDK) Added to the GameLift Plugin")
        else:
            logging.error(f"(BUILDING SDK) Issue building the GameLift Server SDK do to file missing: {GAMELIFT_SERVER_SDK_PATH_LINUX}")
            sys.exit(1)
            
def unzip_files():
    """
    This function will unzip the files we downloaded from Amazon GameLift
    """
    # Unzip the main package
    with zipfile.ZipFile(GAMELIFT_PLUGIN_DOWNLOAD_PATH, 'r') as zip_ref:
        zip_ref.extractall(GAMELIFT_PLUGIN_EXTRACT_PATH)

    # Unzip the main plugin
    with zipfile.ZipFile(GAMELIFT_PLUGIN_MAIN_ZIP_PATH, 'r') as zip_ref:
        zip_ref.extractall(GAMELIFT_PLUGIN_MAIN_PATH)
    
    # We need to download this as its not able to un-zip correctly in Linux with the Release .ZIP
    # However the PPD download link unzips correctly on Linux.
    
    # Download the GameLift-Cpp-ServerSDK-5.1.2.zip file
    # Lets make a new folder
    GAMELIFT_SERVER_SDK_MAIN_PATH.mkdir(parents=True, exist_ok=True)
    
    download_file(GAMELIFT_CCP_SERVER_SDK_ZIP_URL, GAMELIFT_CCP_SERVER_SDK_ZIP_PATH)
    
    with zipfile.ZipFile(GAMELIFT_CCP_SERVER_SDK_ZIP_PATH, 'r') as zip_ref:
        zip_ref.extractall(GAMELIFT_SERVER_SDK_MAIN_PATH)
        
    logging.info("(CREATING PLUGIN) Done! Now Let's check to see if you have all the SDK build requirments:")


def main():
    """
    This function is the main function of the script.
    Lets start by downloading the GameLift Plugin .zip release
    target_operating_system will be Win64, Linux, or OSX
    """
    visual_studio_installed_version = 0
    visual_studio_installed_version_name = ""
    open_sss_installed = None
    visual_studio_installed = None
    cmake_version = None
    git_installed = None
    
    # Start logging
    logging.info(f"(SCRIPT VERSION): You are using the script version: {VERSION}")
    logging.info(f"(CREATING PLUGIN) Target operating_system: {operating_system}")
    logging.info("(CREATING PLUGIN) Downloading GameLift Plugin .zip release..")
    
    # Download the main files
    download_file(GAMELIFT_PLUGIN_DOWNLOAD_URL, GAMELIFT_PLUGIN_DOWNLOAD_PATH)
    
    # Unzip downloads        
    if Path(GAMELIFT_PLUGIN_DOWNLOAD_PATH).exists():
        
        # Lets unzip the files
        unzip_files()
        
        # Check OpenSSL is installed
        open_sss_installed = check_openssl(operating_system)
        
        if open_sss_installed:
            if operating_system == "Win64":
                # Check if Visual Studio is installed
                logging.info("(CREATING PLUGIN) Let's check if you have Visual Studio installed:\n")
                visual_studio_installed_version = check_for_visual_studio()
                
                if visual_studio_installed_version in VISUAL_STUDIO:
                    visual_studio_installed_version_name = VISUAL_STUDIO[visual_studio_installed_version]
                    logging.info(f"(CREATING PLUGIN) {visual_studio_installed_version_name} is installed")
                    visual_studio_installed = True
                else:
                    logging.info("(CREATING PLUGIN) Visual Studio is NOT installed")
                    visual_studio_installed = False
                
            elif operating_system == "OSX":
                pass
            elif operating_system == "Linux":
                pass
            else:
                logging.info("(CREATING PLUGIN) Visual Studio is NOT installed\n")
                visual_studio_installed = False

            # Check if CMake is installed
            logging.info("(CREATING PLUGIN) Let's check if you have CMake installed:")

            # Lets check the version
            cmake_version = check_cmake_version()
            if cmake_version:
                logging.info(f'(CREATING PLUGIN) CMake version: {subprocess.check_output(["cmake", "--version"])}')
                
                # Check if Git is installed
                logging.info("(CREATING PLUGIN) Let's check if you have Git installed:") 
                git_installed = check_git()
                if git_installed:
                    logging.info(f'(CREATING PLUGIN) Git version: {subprocess.check_output(["git", "--version"])}')
                    
                    # Lets check if all these requirments are met
                    if operating_system == "Win64":
                        if open_sss_installed and visual_studio_installed and cmake_version and git_installed:
                            logging.info("(CREATING PLUGIN) All requirements are met!")
                            ready_to_build_sdk = True
                        else:
                            ready_to_build_sdk = False
                            logging.error("(CREATING PLUGIN) All requirements are NOT met.")
                            logging.error("(CREATING PLUGIN) Please install all the requirements and try again")
                    elif operating_system == "Linux" or operating_system == "OSX":
                        if open_sss_installed and cmake_version and git_installed:
                            logging.info("(CREATING PLUGIN) All requirements are met!")
                            ready_to_build_sdk = True
                        else:
                            ready_to_build_sdk = False
                            logging.error("(CREATING PLUGIN) All requirements are NOT met.")
                            logging.error("(CREATING PLUGIN) Please install all the requirements and try again.") 
                    if ready_to_build_sdk:
                        build_sdk(visual_studio_installed_version_name, operating_system)
                        # We are done!
                        logging.info("\n(CREATING PLUGIN) ----------Done! All Done!")
                        # Lets show them the path to the directory
                        logging.info(f"(CREATING PLUGIN) The path to the directory is: {GAMELIFT_PLUGIN_EXTRACT_PATH}")
                    else:
                        logging.error("(CREATING PLUGIN) NOT ready to build SDK!")
                        sys.exit(1)
                else:
                    logging.error("(CREATING PLUGIN) Git is NOT installed! Install GIT from https://git-scm.com/downloads/")
                    webbrowser.open_new_tab(GIT_URL)
                    sys.exit(1)
            else:
                logging.error("(CREATING PLUGIN) CMake is NOT installed! Install CMAKE from https://cmake.org/download/. When installing please add cmake to path.")
                # Open URL
                webbrowser.open_new_tab(CMAKE_URL)
                sys.exit(1)
        else:
            logging.error("(CREATING PLUGIN) OpenSSL is NOT installed! Install the full version of OpenSSL from https://wiki.openssl.org/index.php/Binaries for the appropriate operating_system.")
            webbrowser.open_new_tab(OPEN_SSL_URL)
            sys.exit(1)
    else:
        logging.error("(CREATING PLUGIN) Issue downloading the GameLift Plugin .zip release. You can find the official Amazon GameLift documentation [here](https://aws.amazon.com/documentation/gamelift/).")
        webbrowser.open_new_tab(GL_DOCS_URL)
        sys.exit(1)

if __name__ == "__main__":
    main()

    