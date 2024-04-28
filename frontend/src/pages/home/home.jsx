import { useState } from "react";
import { Link } from "react-router-dom";

export default function Home() {
  const [inputValue, setInputValue] = useState("");

  const handleChange = (event) => {
    setInputValue(event.target.value);
  };

  const handleSubmit = (event) => {
    event.preventDefault();
    fetch("http://localhost:8080/submit", {
      method: "POST",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body: `inputValue=${encodeURIComponent(inputValue)}`,
    })
      .then((response) => response.text())
      .then((data) => {
        console.log(data); // Log response from Go server
      })
      .catch((error) => {
        console.error("Error:", error);
      });
  };
  const handleClear = () => {
    setInputValue(""); // Clear input value
  };

  return (
    <div style={{ display: "flex" }}>
      <nav>
        <ul>
          <li>
            <Link to="/">Home</Link>
          </li>
          <li>
            <Link to="/DiskCreen">Commands</Link>
          </li>
          <li>
            <Link to="/page2">Page 2</Link>
          </li>
        </ul>
      </nav>
      <div>
        <p>Hola mundo</p>
        <br />
        <form onSubmit={handleSubmit}>
          <input
            type="text"
            value={inputValue}
            onChange={handleChange}
            placeholder="Enter something"
          />
          <button type="submit">Submit</button>
          <button type="button" onClick={handleClear}>Clear</button>
        </form>
      </div>
    </div>
  );
}
