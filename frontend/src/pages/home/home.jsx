
import { useState } from "react";
import { Link } from "react-router-dom";
import 'bootstrap/dist/css/bootstrap.min.css';

export default function Home() {
  const [inputValue, setInputValue] = useState("");

  const handleChange = (event) => {
    setInputValue(event.target.value);
  };

  const handleSubmit = (event) => {
    event.preventDefault();
    fetch("http://localhost:8081/submit", {
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
    <div>
      <nav className="navbar navbar-expand-lg navbar-light bg-light">
        <div className="container-fluid">
          <Link className="navbar-brand" to="/">PANTALLA1</Link>
          <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
            <span className="navbar-toggler-icon"></span>
          </button>
          <div className="collapse navbar-collapse" id="navbarSupportedContent">
            <ul className="navbar-nav me-auto mb-2 mb-lg-0">
              <li className="nav-item">
                <Link className="nav-link" to="/DiskCreen">PANTALA2</Link>
              </li>
              <li className="nav-item">
                <Link className="nav-link" to="/Reportes">PANTALLA3</Link>
              </li>
            </ul>
          </div>
        </div>
      </nav>
      <div style={{ padding: '20px' }}>
        <form onSubmit={handleSubmit}>
          <textarea
            className="form-control"
            value={inputValue}
            onChange={handleChange}
            placeholder="Enter something"
            rows={5} // Set the number of visible rows
          />
          <div className="mt-3">
            <button className="btn btn-outline-success me-2" type="submit">Submit</button>
            <button className="btn btn-outline-danger" type="button" onClick={handleClear}>Clear</button>
          </div>
        </form>
      </div>
    </div>
  );
}
