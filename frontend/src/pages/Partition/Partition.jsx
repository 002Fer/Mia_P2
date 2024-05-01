import partitionIMG from "../../assets/particion.png";
import { useState, useEffect } from "react";
import { Link, useParams, useNavigate } from "react-router-dom";
import { useLocation } from "react-router-dom";
export default function Partition({ ip = "localhost" }) {
  const { id } = useParams()
  const [data, setData] = useState([]) 
    const navigate = useNavigate()

    
    // execute the fetch command only once and when the component is loaded
    useState(() => {
 

      var dataF = {
        User: 'root',
        Password: 'root'
      }
      
      console.log(`fech to http://${ip}:8081/`)
      fetch(`http://${ip}:8081/particiones`, {
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
      console.log("click", objIterable)
      navigate(`/login/${id}/${objIterable}`)
    }
    const location = useLocation();

return (
  <>
    
    <p>{location.pathname}</p>
    <br />
    <Link to="/DiskCreen">Commands</Link>
    <br />
    <br />
    <br />
    <br />


    <div style={{ border: "red 1px solid", display: "flex", flexDirection: "row" }}>

      {
        data.map((objIterable, index) => {
          return (
            <div key={index} style={{
              border: "green 1px solid",
              display: "flex",
              flexDirection: "column", // Alinea los elementos en columnas
              alignItems: "center", // Centra verticalmente los elementos
              maxWidth: "100px",
            }}
              onClick={() => onClick(objIterable)}
            >
              <img src={partitionIMG} alt="disk" style={{ width: "100px" }} />
              <p>{objIterable}</p>
            </div>
          )
        })
      }

    </div>
  </>
)
}