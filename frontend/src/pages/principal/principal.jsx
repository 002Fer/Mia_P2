import 'bootstrap/dist/css/bootstrap.min.css';
import diskIMG from "../../assets/img.png";
import { Link } from "react-router-dom";

export default function MyComponent() {
  const handleOpenImage = () => {
    window.open("ruta_de_otra_imagen.jpg", "_blank");
  };

  return (
    <>
      <p>Hola mundo</p>
      <br/>
      <Link to="/DiskCreen">Commands</Link>

      <div className="container">
        <div className="row">
          <div className="col">
            <img src={diskIMG} className="img-fluid img-thumbnail" alt="img" style={{ maxWidth: "200px", maxHeight: "200px" }} />
            <br />
            <button className="btn btn-primary mt-3" onClick={handleOpenImage}>Abrir otra imagen</button>
          </div>
        </div>
      </div>
    </>
  );
}
