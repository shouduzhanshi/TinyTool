#!/bin/bash

export PROJECT_DIR=$(cd `dirname $0`; pwd)

export ANDROID_DIR=../../AndroidStudioProjects/Elephant

export ANDROID_BUILD_CONFIG="$PROJECT_DIR/app.json"

export APPLICATION_ID=$(echo $(cat "$ANDROID_BUILD_CONFIG" | jq ".build.applicationId") | sed 's/"//g')

adb shell pm path "$APPLICATION_ID"

if [ $? = 0 ]; then

    CHANGE_FILE="$1"

    FILE_NAME=${CHANGE_FILE##*/}

    adb shell mkdir sdcard/Android/data/"$APPLICATION_ID"

    adb shell mkdir sdcard/Android/data/"$APPLICATION_ID"/cache

    adb push "$1" sdcard/Android/data/"$APPLICATION_ID"/cache

    adb shell am start -n "$APPLICATION_ID"/com.sunmi.android.elephant.api.container.ContainerActivity --es hotReLoad "hotReLoad://cache/$FILE_NAME" --es router

  else

    cd "$ANDROID_DIR"

    ./gradlew assembleDebug

    adb install -r build/outputs/apk/debug/app-debug.apk

    if [ $? = 0 ]; then

      adb shell am start -n "$APPLICATION_ID"/com.sunmi.android.elephant.core.splash.SplashActivity

    fi

fi
