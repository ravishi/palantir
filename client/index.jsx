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
        const socket = new Socket("ws://localhost:8080/ws", {logger: console.log.bind(console)});

        socket.connect({greetings: 'mellon'});

        this.channel = socket.channel('room:lobby');

        this.channel.on("new_msg", this.onNewMessage.bind(this));

        this.channel.join()
            .receive("ok", (response) => this.setState({messages: (response || {}).messages || []}))
            .receive("error", ({reason}) => console.log("failed to join chat room", reason))
            .receive("timeout", () => console.log("failed to join chat room: timeout"));
    }

    onNewMessage(msg) {
        console.log('message received');
        this.setState({
            messages: [...this.state.messages, msg.Body]
        })
    }

    sendMessage(msg) {
        this.channel.push("new_msg", {Body: msg}, 3000)
            .receive("ok", () => console.log('message sent'))
            .receive("error", (reasons) => console.error("failed to send message", reasons))
            .receive("timeout", () => console.error("failed to send message: timeout"))
    }

    onSubmitMessage() {
        this.sendMessage(this.input.value);
    }

    render() {
        const {messages} = this.state;

        return (
            <div>
                <input type="text" ref={ref => this.input = ref } />
                <button type="button" onClick={this.onSubmitMessage.bind(this, 'sent')}>send</button>
                <div>
                    {<pre>{JSON.stringify(messages, null, 2)}</pre>}
                </div>
            </div>
        )
    }
}

ReactDOM.render(<Chat />, document.getElementById('root'));