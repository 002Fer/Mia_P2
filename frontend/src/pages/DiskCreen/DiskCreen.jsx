
import diskIMG from "../../assets/disk.png";
import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";



  export default function DiskCreen({ip="localhost"}) {
    const [data, setData] = useState([]) 
    const navigate = useNavigate()
    const [inputValue, setInputValue] = useState("");
    
    // execute the fetch command only once and when the component is loaded
    useState(() => {
 

      var dataF = {
        User: 'root',
        Password: 'root'
      }
      
      console.log(`fech to http://${ip}:8081/`)
      fetch(`http://${ip}:8081/tasks`, {
        method: 'POST', 
        headers: {
          'Content-Type': 'application/json' 
        },
        body: JSON.stringify(dataF)
      })
      .then(response => response.json())
      .then(rowdata => {
        console.log(rowdata); // Do something with the response
        setData(rowdata.List)
      })
      .catch(error => {
        console.error('There was an error with the fetch operation:', error);
      });
    }, [])
  const onClick = (objIterable) => {
    //e.preventDefault()
    console.log("click",objIterable)
    navigate(`/disk/${objIterable}`)
    
    fetch("http://localhost:8081/nom_disco", {
      method: "POST",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body: `inputValue=${encodeURIComponent(objIterable)}`,
    })
      .then((response) => response.text())
      .then((data) => {
        console.log(data); // Log response from Go server
      })
      .catch((error) => {
        console.error("Error:", error);
      });
  };
  
  return (
    <>
      <nav className="navbar navbar-expand-lg navbar-light bg-light">
        <div className="container-fluid">
          <Link className="navbar-brand" to="/">PANTALLA1</Link>
          <div className="collapse navbar-collapse" id="navbarSupportedContent">
            <ul className="navbar-nav me-auto mb-2 mb-lg-0">
              <li className="nav-item">
                <Link className="nav-link" to="/DiskCreen">PANTALLA2</Link>
              </li>
              <li className="nav-item">
                <Link className="nav-link" to="/Reportes">PANTALLA3</Link>
              </li>
            </ul>
          </div>
        </div>
      </nav>
      
      <div style={{ backgroundColor: "black", padding: "20px" }}>
        <div className="row">
          {data.map((objIterable, index) => (
            <div className="col" key={index} onClick={() => onClick(objIterable)}>
              <div style={{ border: "1px solid green", 
              display: "flex", 
              flexDirection: "column", 
              alignItems: "center", 
              padding: "10px" 
              }}>
                <img src={diskIMG} alt="disk" style={{ width: "100px" }} />
                
                <p>{objIterable}</p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </>
  );
}