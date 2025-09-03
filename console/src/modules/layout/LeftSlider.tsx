import {Button, Space, Tooltip} from "react-windy-ui";
import HomeIcon from "@/modules/icon/HomeIcon";
import LeftMenu from "@/modules/layout/LeftMenu";
import {useTranslation} from "react-i18next";
import {SiGooglecontaineroptimizedos} from "react-icons/si";
import MenuData from "@/modules/common/MenuData.ts";
import {MenuItem} from "@/types";
import {useAtom} from "jotai";
import {selectedModule} from "@/modules/state/global.ts";
import {useNavigate} from "react-router";

export default function LeftSlider(props) {
    const {collapse} = props;
    const navigate = useNavigate();
    const {t} = useTranslation();
    // const { mdWindow } = useContext(WindowContext);
    const [module, setModule] = useAtom(selectedModule)

    const chooseModule = (moduleName: string) => {
        setModule(moduleName);
        navigate(`/home/${moduleName}/index`);
    }

    const menuData = MenuData()

    return (
        <div className="c-slider" style={{left: '0rem', display: 'flex', flexDirection: "row"}}>
            <div className="c-slider-catalog">
                {
                    menuData.sort((a: MenuItem, b: MenuItem) => a.order - b.order).map(menu => {
                        return <Tooltip key={menu.name} body={menu.displayName} position="right">
                            <Button hasBox={false} inverted hasRipple active={menu.name == module}
                                    onClick={() => chooseModule(menu.name)}
                                    leftIcon={<SiGooglecontaineroptimizedos className="custom-icon"/>} color="purple">
                            </Button>
                        </Tooltip>
                    })
                }

            </div>
            <div className="c-slider-content" style={{flexFlow: '1 1 auto'}}>
                <div className="c-slider-inner">
                    <div className="slider-title">
                        <Space>
                            <HomeIcon/>
                            {!collapse && <span style={{whiteSpace: 'nowrap'}}>My 工具</span>}
                        </Space>
                    </div>
                    <LeftMenu collapse={collapse}/>
                </div>
            </div>
        </div>
    );
}
