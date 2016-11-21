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

class UserRegistrationForm extends React.Component {
	constructor(props) {
		super(props);
		this.state = {name: '', email: ''};
	}
	handleNameChange(event) {
		this.setState({name: event.target.value});
	}
	handleEmailChange(event) {
		this.setState({email: event.target.value});
	}
	handleSubmit(event) {
		event.preventDefault()
		var self = this;
		var user = {
			name: this.state.name,
			email: this.state.email,
		}
		fetch(new Request("/api/user"), {method: 'POST', body: JSON.stringify(user)}).then((resp) => {
			if (resp.ok) {
				self.props.userRegistered()
			}
		});
	}
	render() {
		return e(
			Form,
			{
				inline: true,
				onSubmit: this.handleSubmit.bind(this)
			},
			[
				e(FormControl, {key: 'name', type: 'password', onChange: this.handleNameChange.bind(this)}),
				e(FormControl, {key: 'email', type: 'password', onChange: this.handleEmailChange.bind(this)}),
				e(Button, {key: 'submit', type: 'submit'}, 'Submit')
			]
		);
	}
}

class LoginForm extends React.Component {
	constructor(props) {
		super(props);
		this.state = {token: ''};
	}
	handleTokenChange(event) {
		this.setState({token: event.target.value});
	}
	handleSubmit(event) {
		event.preventDefault()
		var self = this;
		fetch(new Request("/api/sessions"), {method: 'POST', body: this.state.token}).then((resp) => {
			if (resp.ok) {
				self.props.sessionCreated()
			}
		});
	}
	render() {
		return e(
			Form,
			{
				inline: true,
				onSubmit: this.handleSubmit.bind(this)
			},
			[
				e(FormControl, {key: 'token', type: 'password', onChange: this.handleTokenChange.bind(this)}),
				e(Button, {key: 'submit', type: 'submit'}, 'Submit')
			]
		);
	}
}

class Spinner extends React.Component {
	render() {
		return e('p', {}, 'Loading...')
	}
}

class App extends React.Component {
	constructor(props) {
		super(props);
		this.state = {
			waiting: true,
			session: false,
			user: null
		};
	}
	fetchUser() {
		var self = this;
		fetch(new Request("/api/user")).then((resp) => {
			if (resp.ok) {
				self.setState({user: resp.body.json(), session: true, waiting: false});
			}
			if (resp.status === 404) {
				self.setState({session: true, waiting: false});
			}
			if (resp.status == 401) {
				self.setState({waiting: false});
			}
		});
	}
	componentDidMount() {
		this.fetchUser()
	}
	render() {
		if (this.state.waiting) {
			return e(Spinner)
		}
		if (this.state.user) {
			return e('div', {}, [
				e('h1', {}, 'Hello ' + this.state.user),
				e(OrderForm)
			])
		}
		if (this.state.session) {
			return e('div', {}, [
				e(UserRegistrationForm, {userRegistered: this.fetchUser.bind(this)})
			])
		}
		return e(LoginForm, {sessionCreated: this.fetchUser.bind(this)})
	}
}

ReactDOM.render(e(OrderForm), document.getElementById('app'));
