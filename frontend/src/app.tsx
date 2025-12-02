import './App.css'
import logo from "./assets/images/logo-universal.png"
import {Greet, GetAccountBalance} from "../wailsjs/go/main/App";
import {useState} from "preact/hooks";
import {h} from 'preact';

export function App(props: any) {
    const [resultText, setResultText] = useState("Please enter your name below 👇");
    const [name, setName] = useState('');
    const [balance, setBalance] = useState<number | null>(null);

    const updateName = (e: any) => setName(e.target.value);
    const updateResultText = (result: string) => setResultText(result);

    function greet() {
        Greet(name).then(updateResultText);
    }

    function checkBalance() {
        GetAccountBalance().then((result) => setBalance(result)).catch((err) => console.error(err));
    }

    return (
        <>
            <div id="App">
                <img src={logo} id="logo" alt="logo"/>
                <div id="result" className="result">{resultText}</div>
                <div id="input" className="input-box">
                    <input id="name" className="input" onChange={updateName} autoComplete="off" name="input"
                           type="text"/>
                    <button className="btn" onClick={greet}>Greet</button>
                </div>

                <div className="balance-box" style={{ marginTop: '20px' }}>
                    <button className="btn" onClick={checkBalance}>Check Balance</button>
                    {balance !== null && <div className="result">Equity: ${balance.toFixed(2)}</div>}
                </div>
            </div>
        </>
    )
}
