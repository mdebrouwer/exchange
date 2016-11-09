import './app.less'
import React from 'react'
import ReactDOM from 'react-dom'
import { Button, Form, FormControl } from 'react-bootstrap';

const e = React.createElement;

class OrderForm extends React.Component {
	constructor(props) {
		super(props);
		this.state = {details: ''};
	}
	handleDetailsChange(event) {
		this.setState({details: event.target.value});
	}
	handleSubmit(event) {
		event.preventDefault()
		fetch(new Request("/api/orders"), {method: 'POST', body: JSON.stringify({order: this.state.details})}).then((resp) => { console.log(resp) });
	}
	render() {
		return e(
			Form,
			{
				inline: true,
				onSubmit: this.handleSubmit.bind(this)
			},
			[
				e(FormControl, {key: 'details', type: 'text', placeholder: 'order details', onChange: this.handleDetailsChange.bind(this)}),
				e(Button, {key: 'submit', type: 'submit'}, 'Submit')
			]
		);
	}
}

ReactDOM.render(e(OrderForm), document.getElementById('app'));
