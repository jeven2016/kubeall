import {Box, Pagination, Table, Typography} from "react-windy-ui";
import {useCallback} from "react";
import {useLocation, useNavigate} from "react-router";
import {ListParam, PageList} from "@/types";


const listSites = () => {
    return {payload: [], code: "200", message: "nothing"}
};


export default function ResourceList({listParam, pageList}: { listParam: ListParam, pageList: PageList }) {
    const navigate = useNavigate();
    const location = useLocation();


    //
    // if (isSuccess && data?.code != "200") {
    //     Notification.error({
    //         body: data?.message,
    //         position: "topCenter"
    //     });
    // }

    const enter = useCallback((id: string) => {
        navigate(`/home/sites/${id}`);
    }, []);


    return (
        <div className="c-content-area">
            <div className="c-content-header">
                <div className="c-header-panel">
                    <div className="c-book-chanel"> {listParam?.title}</div>
                    <div className="c-header-desc">
                        <Typography italic>{listParam.description}</Typography>
                    </div>
                </div>
            </div>

            <div className={"c-content-body"}>
                <div>
                    <Box block={true}
                         left={listParam.buttons}
                         right={listParam.filterArea && <div>
                             {listParam.filterArea}
                         </div>}
                    />

                </div>
                <div className="c-book-chanel-list">
                    <div className="c-table-content">
                        <Table
                            type="striped"
                            loadData={listParam.data}
                            cells={listParam.cells}
                            checkable={true}
                            checkType={"checkbox"}
                            checkedRows={listParam.checkedRows}
                            onCheckAll={listParam.onCheckAll}
                            onCheckChange={listParam.onCheck}
                            hover={true}
                            hasBorder={true}
                        />
                    </div>
                    <div>
                        {pageList &&
                            <Pagination
                                pageCount={pageList.totalPages ?? 0}
                                page={pageList.page}
                                pageRange={pageList.pageSize}
                                siblingCount={1}
                                hasPageRange={true}
                                pageRanges={[10, 20, 50]}
                            />
                        }
                    </div>
                </div>
            </div>
        </div>
    );
}

