import QtQuick 2.7
import Ubuntu.Components 1.3
//import QtQuick.Controls 2.2
import QtQuick.Layouts 1.3
import Qt.labs.settings 1.0
//import "Components"

MainView {
    id: root
    objectName: 'mainView'
    applicationName: 'whats-ut.walking-octopus'
    automaticOrientation: true

    width: units.gu(80)
    height: units.gu(60)

    PageStack {
        id: mainStack
        Component.onCompleted: mainStack.push(Qt.resolvedUrl("Login.qml"))
    }
    //Button {
        //text: "Manual reload"
        //color: UbuntuColors.green
        //onClicked: internal.refreshPage()
        //Layout.fillWidth: true
    //}
}
