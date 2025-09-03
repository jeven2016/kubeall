import {Menu} from 'react-windy-ui';
import IconSystem from '@/modules/icon/SystemIcon';
import {useTranslation} from 'react-i18next';
import {useCallback, useMemo} from "react";
import {MenuItemIdType} from "react-windy-ui/dist/menu/Menu";
import {selectedMenuItem, selectedModule} from "@/modules/state/global.ts";
import {useAtom} from "jotai";
import MenuData from "@/modules/common/MenuData.ts";
import {useNavigate} from "react-router";

const menus = [
    {name: "下载工具", id: "/home/download-tool"},
    {name: "系统设置", id: "/home/settings"},
    {name: "其他工具", id: "/home/other"},
]

export default function LeftMenu(props) {
    const {collapse} = props;
    const {t} = useTranslation();
    const [module, switchModule] = useAtom(selectedModule);
    const [selectedItem, setSelectedItem] = useAtom(selectedMenuItem);
    const menuData = MenuData()
    const navigate = useNavigate();
    const menus = useMemo(() => {
        const m = menuData.filter(menu => menu.name == module);
        if (m.length > 0) {
            return m[0].children || []
        }
        return []
    }, [module, menuData])

    const selectItem = useCallback((id: MenuItemIdType) => {
        setSelectedItem(id as string)
        navigate(`/home/${module}/${id}`)
    }, [setSelectedItem, module])

    return (
        <Menu
            activeItems={selectedItem}
            onSelect={selectItem}
            compact={collapse}
            defaultActiveItems={['item1']}
            hasBox={false}
            type="primary"
            hasRipple={true}>
            {
                menus.map((menu) => {
                    const chr = menu.children || [];
                    if (chr.length == 0) {
                        return <Menu.Item id={menu.name} key={menu.name} icon={<IconSystem/>}>
                            {menu.displayName}
                        </Menu.Item>
                    }
                    return <Menu.SubMenu header={menu.displayName} id={menu.name}>
                        {
                            chr.map((item) => (
                                <Menu.Item id={item.name} key={item.name} icon={<IconSystem/>}>
                                    {item.displayName}
                                </Menu.Item>
                            ))
                        }

                    </Menu.SubMenu>

                })
            }
        </Menu>
    );
}
