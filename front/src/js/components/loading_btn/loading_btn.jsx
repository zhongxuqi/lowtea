import React from 'react'

import Language from '../../language/language.jsx'

import './loading_btn.less'

export default class LoadingBtn extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            status: "active",
        }
    }

    button(status) {
        this.setState({status:status})
    }

    render() {
        return (
            <div className="lowtea-loading-btn" style={{display:{true:"none", false:"block"}[this.state.status=="finish"]}}>
                <button type="button" className="btn btn-default" onClick={(()=>{
                    this.setState({status:"loading"})
                    this.props.onClick()
                }).bind(this)}><i className="fa fa-spinner fa-spin fa-lg fa-fw" style={{display:{true:"inline-block", false:"none"}[this.state.status=="loading"]}}></i>{{true:Language.textMap("Loading..."), false:Language.textMap("Load More")}[this.state.status=="loading"]}</button>
            </div>
        )
    }
}
