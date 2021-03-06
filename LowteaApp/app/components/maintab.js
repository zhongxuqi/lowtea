import React, {Component} from 'react'
import {
    StyleSheet,
    View,
    Text,
    TouchableHighlight,
} from 'react-native'
import Icon from 'react-native-vector-icons/FontAwesome';
import BaseCSS from '../config/css.js'

export default class MainTab extends Component {
    constructor(props) {
        super(props)
        this.state = {
            press: false,
        }
    }

    render() {
        return (
            <TouchableHighlight style={BaseCSS.container} underlayColor={BaseCSS.colors.transparent} 
                onHideUnderlay={(()=>{
                    this.setState({press: false})
                }).bind(this)}    
                onShowUnderlay={(()=>{
                    this.setState({press: true})
                }).bind(this)}    
                onPress={this.props.onPress}>
                <View style={{false:styles.button,true:styles.button_active}[this.props.active||this.state.press]}>
                    <Icon name={{false:this.props.icon,true:this.props.active_icon}[this.props.active||this.state.press]}
                        size={20} color={{false:BaseCSS.colors.black,true:BaseCSS.colors.blue}[this.props.active||this.state.press]}/>
                    <Text style={{false:styles.text,true:styles.text_active}[this.props.active||this.state.press]}>{this.props.title}</Text>
                </View>
            </TouchableHighlight>
        )
    }
}

const styles=StyleSheet.create({
    button: {
        backgroundColor: BaseCSS.colors.transparent,
        borderTopColor: BaseCSS.colors.separation_line,
        borderTopWidth: 1,
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center',
        paddingVertical: 7,
    },
    button_active: {
        backgroundColor: BaseCSS.colors.transparent,
        borderTopColor: BaseCSS.colors.separation_line,
        borderTopWidth: 1,
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center',
        paddingVertical: 7,
    },
    text: {
        color: BaseCSS.colors.black,
    },
    text_active: {
        color: BaseCSS.colors.blue,
    },
})
