import { useState } from "react";

export default function SignIn({ ip = "localhost" }) {
   const [error, setError] = useState(null);

   const handleSubmit = async (e) => {
      e.preventDefault();
      
      const user = e.target.uname.value;
      const pass = e.target.psw.value;

      try {
         const response = await fetch(`http://${ip}:8081/login`, {
            method: 'POST',
            headers: {
               'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username: user, password: pass })
         });

         if (!response.ok) {
            throw new Error('Failed to login');
         }

         // Redirigir al usuario a otra página si el inicio de sesión es exitoso
         window.location.href = '/dashboard'; // Cambia '/dashboard' por la ruta a la que quieres redirigir
      } catch (error) {
         console.error('Error:', error);
         setError('Invalid username or password');
      }
   }

   return (
      <div className="container mt-5">
         <h2>Login</h2>

         <form onSubmit={handleSubmit}>
            <div className="mb-3">
               <label htmlFor="uname" className="form-label">Username</label>
               <input type="text" className="form-control" id="uname" placeholder="Enter Username" name="uname" required/>
            </div>

            <div className="mb-3">
               <label htmlFor="psw" className="form-label">Password</label>
               <input type="password" className="form-control" id="psw" placeholder="Enter Password" name="psw" required/>
            </div>
               
            {error && <div style={{ color: 'red' }}>{error}</div>}

            <button type="submit" className="btn btn-primary">Login</button>
         </form>
      </div>
   )
}
