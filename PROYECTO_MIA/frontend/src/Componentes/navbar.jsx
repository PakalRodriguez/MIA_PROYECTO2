import React, { useState } from 'react';
import hacker from '../iconos/persona.png';
import Comandos from '../Paginas/comandos';
import Explorer from '../Paginas/discos';

export default function Navbar(){
    const [componenteActivo, setComponenteActivo] = useState(<Comandos/>);

    function comandos(idPagina){
        let componente
        if (idPagina === 1){
            componente =  <Comandos/>
        }else if (idPagina === 2){
            componente = <Explorer/>
        }
        setComponenteActivo(componente)
    }

    return(
        <>
            <nav className="navbar navbar-expand-lg navbar-dark bg-dark">
                {/*COLUMNAS*/}
                <div id="espacio">&nbsp;&nbsp;&nbsp;</div>
                
                <div className="conteiner-fluid"> 
                    <img src={persona} alt="" width="64" height="64" className="d-inline-block align-text-top"></img>
                </div>

                <div className="conteiner"> 
                    {/*Fila 1 (titulo del proyecto, RESPALDO)*/}
                    <div className="container-fluid">
                    <a className="navbar-brand" href="#">ARCHIVOS PROYECTO 2</a>

                        {/*Cada bloque div aqui dentro es una nueva fila hacia abajo*/}
                        {/*Fila 2 (menus)*/}
                        <div className="collapse navbar-collapse" id="navbarColor02">
                            {/*ul es una lista no ordenada*/}
                            <ul className="navbar-nav me-auto">
                                {/*LISTA DE MENUS QUE ESTARAN EN LA BARRA DE NAVEGACION*/}
                                <li className="nav-item">
                                    <a className="nav-link active" type="button" onClick={() => comandos(1)}>Comandos</a>
                                </li>

                                <li className="nav-item">
                                    <a className="nav-link" type="button" onClick={() => comandos(2)}>Explorador</a>
                                </li>

                                <li className="nav-item">
                                    <a className="nav-link" type="submit">Logout</a>
                                </li>

                            </ul>{/*Fin de lista de menus*/}
                        </div>{/*Fila de menus en la barra de navegacion*/}
                    </div>{/*Fila Titulo*/}
                </div>{/*Cierro tercer columna (Menu)*/}
            </nav> 
            {componenteActivo}
        </>
    );
}
