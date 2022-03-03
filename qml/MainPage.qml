/*
 * Copyright (C) 2022  JohnDoe
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 3.
 *
 * qml-rss-reader is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

import QtQuick 2.9
import Ubuntu.Components 1.3


Page {
    id: viewerPage

    title: "Placeholder"
    header: PageHeader {
        id: header
        title: "Placeholder"
    }

    Label {
        anchors {
            top: header.bottom
            left: parent.left
            right: parent.right
            bottom: parent.bottom
        }
        text: i18n.tr('Successful connection')
//         font.pixelSize: 25

        verticalAlignment: Label.AlignVCenter
        horizontalAlignment: Label.AlignHCenter
    }
}
