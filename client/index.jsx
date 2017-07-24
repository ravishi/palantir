import React from 'react';
import {Socket} from 'phoenix';
import ReactDOM from 'react-dom';

class Chat extends React.Component {
    constructor(...args) {
        super(...args);
        this.state = {
            messages: []
        };
    }

    componentDidMount() {
        const socket = new Socket("ws://localhost:8080/ws");

        socket.connect({greetings: 'Mellon!'});

        this.channel = socket.channel('room:lobby');

        this.channel.on("new_msg", this.onNewMessage.bind(this));

        this.channel.join()
            .receive("ok", ({messages}) => this.setState({messages}))
            .receive("error", ({reason}) => console.log("failed join", reason))
            .receive("timeout", () => console.log("Networking issue. Still waiting..."));
    }

    onNewMessage(msg) {
        this.setState({
            messages: [...messages, msg]
        })
    }

    sendMessage(msg) {
        this.channel.push("new_msg", {body: msg}, 3000)
            .receive("ok", this.onNewMessage.bind(this))
            .receive("error", (reasons) => console.log("create failed", reasons))
            .receive("timeout", () => console.log("Networking issue..."))
    }

    onSubmitMessage() {
        this.sendMessage(this.input.value);
    }

    render() {
        const {messages} = this.state;

        return (
            <div>
                <input type="text" ref={ref => this.input = ref } />
                <button type="button" onClick={this.onSubmitMessage.bind(this)}>send</button>
                <div>
                    {<pre>{JSON.stringify(messages, '  ')}</pre>}
                </div>
            </div>
        )
    }
}

ReactDOM.render(<Chat />, document.getElementById('root'));