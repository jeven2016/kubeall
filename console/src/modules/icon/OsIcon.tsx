import {Tooltip} from "react-windy-ui";
import {ImWindows} from "react-icons/im";
import React from "react";
import {AiOutlineLinux} from "react-icons/ai";
import {FaUbuntu} from "react-icons/fa";
import {RiCentosFill} from "react-icons/ri";
import {SiRockylinux} from "react-icons/si";

const osIconSize = "1.5rem";
export const osType = {
    windows: "windows",
    linux: "linux",
    ubuntu: "ubuntu",
    centos: "centos",
    kylin: "kylin",
    rocky: "rocky",
}

export default function getOsIcon(os: string|undefined): React.ReactNode {
    if (!os){
        return ""
    }
    let osIcon: React.ReactNode
    let body: string = osType[os]
    os = os.toLowerCase();

    switch (os) {
        case osType.windows:
            osIcon = <ImWindows size={osIconSize} className="os-icon"/>
            break
        case osType.linux:
            osIcon = <AiOutlineLinux size={osIconSize} className="os-icon"/>
            break
        case osType.ubuntu:
            osIcon = <FaUbuntu size={osIconSize} className="os-icon"/>
            break
        case osType.centos:
            osIcon = <RiCentosFill size={osIconSize} className="os-icon"/>
            break
        case osType.rocky:
            osIcon = <SiRockylinux size={osIconSize} className="os-icon"/>
            break
    }
    return <Tooltip body={body} position="bottom">
        <span>{osIcon} </span>
    </Tooltip>
}