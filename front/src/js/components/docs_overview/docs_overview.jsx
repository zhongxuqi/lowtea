import React from 'react'

import SearchBar from '../searchbar/searchbar.jsx'
import DocsList from '../docs_list/docs_list.jsx'
import TopFlagList from '../topflag_list/topflag_list.jsx'
import TopStarList from '../topstar_list/topstar_list.jsx'
import LoadingBtn from '../loading_btn/loading_btn.jsx'
import DocListTitle from '../doc_list_title/doc_list_title.jsx'

import HttpUtils from '../../utils/http.jsx'
import Language from '../../language/language.jsx'

import './docs_overview.less'

export default class DocsOverView extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            documents: [],
            pageSize: 10,
            pageIndex: 0,
            docTotal: 0,
            keyword: "",
            topFlagDocuments: [],
            topStarDocuments: [],
        }
        this.getDocuments(this.state.pageSize, this.state.pageIndex)

        HttpUtils.get("/api/member/top_star_documents",{},((resp)=>{
            let topStarDocuments = []
            for (let i=0;i<resp.documents.length;i++) {
                let document = resp.documents[i]
                document.starPercent = Math.round(100 * document.starNum / resp.memberNum)
                topStarDocuments.push(document)
            }
            this.setState({topStarDocuments:topStarDocuments})
        }).bind(this), (resp)=>{
            HttpUtils.alert("["+resp.status+"] "+resp.responseText)
        })
        
        HttpUtils.get("/api/member/top_flag_documents",{},((resp)=>{
            let topFlagDocuments = []
            for (let i=0;i<resp.documents.length;i++) {
                let document = resp.documents[i]
                document.flagPercent = Math.round(100 * document.flagNum / resp.adminNum)
                topFlagDocuments.push(document)
            }
            this.setState({topFlagDocuments:topFlagDocuments})
        }).bind(this), (resp)=>{
            HttpUtils.alert("["+resp.status+"] "+resp.responseText)
        })
    }

    getDocuments(pageSize, pageIndex, successFunc, errorFunc) {
        HttpUtils.get("/api/member/documents", {
            pageSize: pageSize,
            pageIndex: pageIndex,
            keyword: this.state.keyword,
        }, ((resp) => {
            if (resp.documents == null) resp.documents = []
            if (pageIndex == 0) {
                this.setState({
                    documents: resp.documents,
                    docTotal: resp.docTotal,
                })
            } else {
                this.setState({
                    documents: this.state.documents.concat(resp.documents),
                    docTotal: resp.docTotal,
                })
            }
            
            if (pageIndex + 1 >= resp.pageTotal) {
                this.refs.LoadingBtn.button("finish")
            } else {
                this.refs.LoadingBtn.button("active")
            }
            this.setState({pageIndex: pageIndex})
        }).bind(this), (resp) => {
            HttpUtils.alert("["+resp.status+"] "+resp.responseText)
            this.refs.LoadingBtn.button("active")
        })
    }

    render() {
        return <div className="lowtea-doc-overview">
            <div className="col-md-9 col-sm-9 col-xs-9 doc-list-container">
                <div className="searchbar-container">
                    <SearchBar ref="searchbar" onClick={(()=>{
                        this.state.keyword = this.refs.searchbar.getValue()
                        this.state.pageIndex = 0
                        this.getDocuments(this.state.pageSize, 0)
                        this.setState({})
                    }).bind(this)}></SearchBar>
                </div>
                
                <div className="clearfix" style={{margin:"0px 30px", paddingBottom:"10px"}}>
                    <DocListTitle docTotal={this.state.docTotal}></DocListTitle>

                    <DocsList documents={this.state.documents}></DocsList>

                    <div style={{width:"100%"}}>
                        <LoadingBtn ref="LoadingBtn" onClick={(()=>{
                            this.getDocuments(this.state.pageSize, this.state.pageIndex + 1) 
                        }).bind(this)}></LoadingBtn>
                    </div>
                </div>
            </div>
            <div className="col-md-3 col-sm-3 col-xs-3" style={{margin:"30px 0px"}}>
                <TopFlagList documents={this.state.topFlagDocuments}></TopFlagList>
                <TopStarList documents={this.state.topStarDocuments}></TopStarList>
            </div>
        </div>
    }
}
