import {IconWarning, Notification} from 'react-windy-ui';

export const Message = {
    error:(title:string, body:string, icon:any)=>{
        const titleInfo:string = title ?? "";
        const bodyInfo:string = body ?? "";

        Notification.simple({
            title: titleInfo,
            body: bodyInfo,
            icon: <IconWarning style={{ color: '#c88f3f' }} />
        })
    },
    ok:(title:string, body:string)=>{
        const titleInfo:string = title ?? "";
        const bodyInfo:string = body ?? "";

        Notification.ok({
            title: titleInfo,
            body: bodyInfo
        })
    }
}