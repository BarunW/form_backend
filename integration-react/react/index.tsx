import ReactDOM from 'react-dom/client';
import React from 'react';
import Form from './Form';
import { PropsType } from './Types';


export function renderIndex(p : PropsType){
    const root = ReactDOM.createRoot(document.getElementById('index') as HTMLElement);
    root.render( 
       <Form {...p}/>
    );
}
