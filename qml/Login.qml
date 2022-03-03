/*
 * Copyright (C) 2022  JohnDoe
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 3.
 *
 * whatsut is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

import QtQuick 2.7
import Ubuntu.Components 1.3
//import QtQuick.Controls 2.2
import QtQuick.Layouts 1.3
import Qt.labs.settings 1.0
import "Components"

Page {
    anchors.fill: parent

    header: PageHeader {
        id: header
        title: i18n.tr('WhatsUT')
    }
    title: "WhatsUT"

    // Todo: Create an adaptive layout
    ColumnLayout {
        anchors {
            top: header.bottom
            left: parent.left
            right: parent.right
            bottom: parent.bottom
        }
        spacing: units.gu(0.5)

        Label {
            text: `Welcome to WhatsUB! It works by imitating WhatsApp Web, so you'll need a second phone or an emulator to log in. You'll need that device only once, and it doesn't have to stay online. To get started, enable Multi-Device Beta and scan this QR code from the official mobile app.`
            wrapMode: "WordWrap"
            Layout.fillWidth: true
            Layout.leftMargin: units.gu(2)
            Layout.rightMargin: units.gu(2)
            Layout.topMargin: units.gu(2)
        }
        //Label {
            //text: qmlBridge.loginToken
            //Layout.fillWidth: true
            //Layout.leftMargin: units.gu(2)
            //Layout.rightMargin: units.gu(2)
            //Layout.topMargin: units.gu(2)
        //}
        UbuntuShape {
            Layout.fillHeight: true
            Layout.fillWidth: true

            Layout.leftMargin: units.gu(3)
            Layout.topMargin: units.gu(3)
            Layout.rightMargin: units.gu(3)
            Layout.bottomMargin: units.gu(3)

            sourceFillMode: UbuntuShape.PreserveAspectFit
            aspect: UbuntuShape.Flat
            backgroundColor: "#ffffff"

            source: QRCode {
                id: qr
                width: parent.width - units.gu(5)
                height: width
                anchors.bottom: parent.bottom
                value: qmlBridge.loginToken?qmlBridge.loginToken:""
            }
        }
        Button {
            text: "Manual reload"
            color: UbuntuColors.green
            onClicked: internal.onTokenChanged()
            Layout.fillWidth: true
        }
    }
    QtObject {
        id: internal

        function onTokenChanged() {
            if (qmlBridge.loginToken == "DONE") {
                mainStack.pop()
                mainStack.push(Qt.resolvedUrl("MainPage.qml"))
            }
        }
    }
    //Timer {
        //interval: 150
        //running: true
        //repeat: true
        //onTriggered: internal.onTokenChanged()
    //}

    Connections {
        target: qmlBridge
        // onLoginTokenChanged: console.info("[whatsut-qml] ", qmlBridge.loginToken)
        onLoginTokenChanged: internal.onTokenChanged()
    }
    Component.onCompleted: console.log("Completed Running!", qmlBridge)
}
