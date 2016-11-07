import './app.less'
import React from 'react'
import ReactDOM from 'react-dom'

const e = React.createElement;

class Greeting extends React.Component {
	render() {
		return e('div', {}, 'Hello!');
	}
}

ReactDOM.render(e(Greeting), document.getElementById('app'));
