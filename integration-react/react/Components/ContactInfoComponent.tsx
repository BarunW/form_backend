import React, { Ref } from "react";
import { ComponentsParams, inputProperties} from "../Types";
import {ReactElement } from "react";
import { InsertAnswer } from "../Form";


export function ContactInfoComponent(params : ComponentsParams, fRef : Ref<HTMLFormElement>) {
    const {questionUUID, questionNumber, question} = params

    const qs= question.setting.question_setting
    const data = question.data
    const handleOnSubmit =(e : any) => {
        e.preventDefault(); 
        const formData =  new FormData(e.currentTarget)
        const answer : Object={
            "first_name" : formData.get("first_name"),
            "last_name" : formData.get("last_name"),
            "contact" : formData.get("contact"),
            "email" : formData.get("email"), 
            "company" : formData.get("company"), 
        }
        // ADD the data to the answers
        InsertAnswer(questionUUID, answer) 
    }

    return (
    <div>
        <div className="flex-col">
        <span>{questionNumber}</span>
        <span>{" -> "}</span>
        <span>{data.question}</span>
        </div>
        <span>{question.data.description}</span>
        
        <form className="flex-col" id="form" ref={fRef} onClick={(e)=>handleOnSubmit(e)}>
            {
                qs[1].include === true ?
                CreateInputElement({label:qs[1].label, type:"text", name: "first_name", placeholder:data["first_name"]}) : null
            }          
            {
                qs[2].include === true ?
                CreateInputElement({label:qs[2].label, type:"text", name:"last_name", placeholder:data["last_name"]}) : null
            }

            {
                qs[3].include === true ?
                CreateInputElement({label:qs[3].label, type:"tel", name:"contact", placeholder:data["phone_number"]["number"]}) : null
            }          
            {
                qs[4].include === true ?
                CreateInputElement({label:qs[4].label, type:"email", name:"email", placeholder:data["email"]}) : null
            }
            {
                qs[5].include == true ?
                CreateInputElement({label:qs[5].label, type:"text", name:"company", placeholder:data["company"]}) : null
            }
            <button type="submit" id="getFormData">OK</button>
        </form>
    </div>
    )
}



function CreateInputElement(prop: inputProperties) : ReactElement{
   return(
        <div className="flex-col">
            <label>{prop.label}</label>
            <input
                name={prop.name}
                type={prop.type}
                placeholder={prop.placeholder}
            /> 
        </div>
   )
}
