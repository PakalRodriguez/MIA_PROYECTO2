//import { useState } from "react";
import { useState, useEffect } from "react";

import texto from '../Iconos/archivo.png';
import carpeta from '../Iconos/carpetas.png';
import volver from '../Iconos/flecha.png';
import "../StyleSheets/explorer.css"

export default function Explorer({newIp="localhost"}){
    const [ archivos, setArchivos ] = useState([]);
    const [ estado, setEstado ] = useState(true); //para evitar que muestre imagen cuando es cocntenido de archivo
    const [ path, setPath ] = useState("path: /");

    useEffect(() => {
        fetch(`http://${newIp}:8080/explorer`)
          .then(response => response.text())
          .then(text => {
            console.log("Texto recibido del backend:", text);
    
            // Validar si contiene el prefijo esperado
            if (text.includes("PARTICIONES MONTADAS:")) {
              const rawList = text.split(":")[1].trim(); // "871A, 871B"
              const ids = rawList.length > 0 ? rawList.split(",").map(id => id.trim()) : [];
              setArchivos(ids);
            } else {
              setArchivos([]);
            }
          })
          .catch(error => {
            console.error('Error en la solicitud Fetch:', error);
            alert("⚠️ No se pudo conectar con el backend o la respuesta fue inválida.");
          });
    }, []);
    

    const onClick = (archivo) => {
        console.log("buscar",archivo)
        let tmp = path+archivo+"/"
        setPath(tmp)
        fetch(`http://${newIp}:8080/contenido`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json'},
            body: JSON.stringify(archivo)
        })
        .then(Response => Response.json())
        .then(rawData => {console.log(rawData); setArchivos(rawData);})
    }

    const getFile = (archivo) => {
        console.log("buscar",archivo)
        let tmp = path+archivo+"/"
        setPath(tmp)
        setEstado(false)
        fetch(`http://${newIp}:8080/file`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json'},
            body: JSON.stringify(archivo)
        })
        .then(Response => Response.json())
        .then(rawData => {console.log(rawData); setArchivos(rawData);})
    }

    const back = () =>{
        let tmp2 = path.split("/")
        if (tmp2.length > 2) {
            let newPath = "path: /"
            for (let i=1; i < tmp2.length-2; i++){
                newPath += tmp2[i]+"/"
            }
            console.log("back ", newPath)
            setPath(newPath)
            setEstado(true) //por si estaba mostrando contenido de un archivo
            fetch(`http://${newIp}:8080/back`)
            .then(Response => Response.json())
            .then(rawData => {console.log(rawData); setArchivos(rawData);})
        }
    }

    return(
        <>
            <div className="container">
                <div className="d-flex justify-content-center">
                    <div className="explorer-card">
                        <div className="explorer-card-header">
                            <img onClick={back} src={volver} alt="volver" style={{width: "20px", margin: "5px"}} />
                            {path}
                        </div>
                        <div className="container-with-scroll" style={{display:"flex", flexDirection:"row"}}>
                            {archivos && archivos.length > 0 ? (
                                archivos.map((archivo, index) => {
                                    return (
                                        estado ? (
                                            <div key={index} style={{
                                                display: "flex",
                                                flexDirection: "column", // Alinea los elementos en columnas
                                                alignItems: "center", // Centra verticalmente los elementos
                                                maxWidth: "100px",
                                                margin: "10px"
                                                }}
                                            >
                                                {archivo.endsWith('.txt')? (
                                                    <img onClick={() => getFile(archivo)} src={texto} alt="archivo" style={{width: "100px"}} />    
                                                ):(
                                                    <img onClick={() => onClick(archivo)} src={carpeta} alt="archivo" style={{width: "100px"}} />
                                                )}
                                                {archivo}
                                            </div>
                                        ):(
                                            <div key={index} style={{
                                                margin:"5px", 
                                                width: "100%", 
                                                maxHeight: "200px", 
                                                wordWrap: "break-word",
                                                overflowY:"auto"
                                            }}>
                                                {archivo}
                                            </div>
                                        ) 
                                    )
                                })
                            ):(
                                <div>No hay archivos disponibles</div>
                            )}
                        </div>
                    </div>
                </div>
            </div>
        </>
    );
}