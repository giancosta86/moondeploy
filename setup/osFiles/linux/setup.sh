#!/bin/bash

#-------------------------
# MOONDEPLOY SETUP SCRIPT
#-------------------------
#
#
# Copyright © Gianluca Costa


set -e


#°°°°°°°°°°°°°°°°°°
# CUSTOM VARIABLES
#°°°°°°°°°°°°°°°°°°

#You can change this value to install MoonDeploy to another directory
TARGET_DIRECTORY="$HOME/MoonDeploy"


#°°°°°°°°°°°°°°°°°°°°°°
# INTERNAL DECLARATIONS
#°°°°°°°°°°°°°°°°°°°°°°

waitYesNo () {
    read userResponse

    case $userResponse in
        [yY])
            true
            ;;
        *)
            false
            ;;
    esac
}

SCRIPT_NAME=$(basename "$0")
SCRIPT_DIR=$(dirname "$0")


#°°°°°°°°°°°
# MAIN FLOW
#°°°°°°°°°°°

reset


#
# Waiting for user input
#

echo "MoonDeploy will be installed to '${TARGET_DIRECTORY}'."
echo
echo
echo -n "Do you wish to continue? [Y/N]: "

waitYesNo

#
# Copying the files
#

mkdir -p "${TARGET_DIRECTORY}"

cp -r "${SCRIPT_DIR}" "${TARGET_DIRECTORY}"
rm "${TARGET_DIRECTORY}/${SCRIPT_NAME}"

chmod 700 "${TARGET_DIRECTORY}/moondeploy"


#
# Registering MIME type
#

mkdir -p "$HOME/.local/share/mime/packages"

cat > "$HOME/.local/share/mime/packages/application-moondeploy.xml" <<- EOM
<?xml version="1.0" encoding="UTF-8"?>
<mime-info xmlns="http://www.freedesktop.org/standards/shared-mime-info">
    <mime-type type="application/moondeploy">
        <comment>MoonDeploy Application</comment>
        <icon name="application-moondeploy"/>
        <glob-deleteall/>
        <glob pattern="*.moondeploy"/>
    </mime-type>
</mime-info>
EOM


#
# Registering application
#

mkdir -p "$HOME/.local/share/applications/"

cat > "$HOME/.local/share/applications/moondeploy.desktop" <<- EOM
[Desktop Entry]
Name=MoonDeploy
Exec=${TARGET_DIRECTORY}/moondeploy
MimeType=application/moondeploy
Icon=${TARGET_DIRECTORY}/moondeploy.png
Terminal=false
Type=Application
Categories=
Comment=Lightweight software deployment
EOM


#
# Updating caches
#

update-desktop-database $HOME/.local/share/applications
update-mime-database $HOME/.local/share/mime


#
# Conclusion
#

echo
echo "Installation successful. You might want to add '${TARGET_DIRECTORY}' to your PATH to run MoonDeploy from within the terminal."
echo
read -p "Press [Enter] key to end..."
echo
echo
