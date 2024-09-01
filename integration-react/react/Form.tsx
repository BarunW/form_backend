import  React, {  Ref, useRef, useState } from 'react';
import { MCQ } from "./Components/MCQComponent"
import './Form.css'
import {  ComponentsParams, PropsType } from './Types';
import { ContactInfoComponent } from './Components/ContactInfoComponent';


let Answers : Object = {}

function componentMapper(questions: any[], fs : Object, formRef : any) : React.ReactElement[]{
    let components : React.ReactElement[] = [] 
    questions.map((quest, idx) => { 
        const datas : ComponentsParams = {
            questionNumber : idx + 1,
            questionUUID   : quest["quuid"],
            question : quest,
            formSettings : fs,
        }

        if (quest["title"] === "Multiple Choice"){
            components.push(MCQ(datas))
        }
        
        if (quest["title"] === "Contact Info"){
            components.push(ContactInfoComponent(datas, formRef))
        }

    })
    return components
}

export function InsertAnswer(qId:Number, answerForQuest: any){
    const key : string = ""+qId
    const exist : boolean = Answers.hasOwnProperty(key)
    if (exist){
        Answers[key]["answer"] = answerForQuest 
    }
}

async function handleSubmitForm(formId:string, accountId : string, formRef : any){
 
     const response = await fetch(`http://${accountId}.localhost:8080/form/${formId}/submissions`, {
         method: 'POST',
         headers: {
             'Content-Type': 'application/json',
             'Access-Control-Allow-Origin': '*',
        },
        body: JSON.stringify(Answers)
    });

    if (!response.ok) {
        throw new Error(`HTTP error: ${response.status}`);
    }

  // Assuming the response is also JSON
    try{
        const responseData = await response.json();
        console.log(responseData)
    }catch(err){
        console.log(err)
    }
    
    if (formRef.current){
        formRef.current.click()
    }
    
}

export default function Form(props:PropsType){
    // Answers Context  
    Answers = props.answers
    console.log(props.questions, props)

    const formRef = useRef<HTMLFormElement>(null)
    const components  = componentMapper(props.questions, props.form_settings, formRef)
    const [count, setCount] = useState(0);    
    const [scrollY, setScrollY] = useState(0);
    
    const [animation, setAnimation] = useState({
        animated : "",
        onAnimation : "animation"
    })

    const [status, setStatus] = useState("")

    function handleOnclickUp(){
        if ( count + 1< components.length){
         setCount(prev => prev + 1);
        }
        setStatus(animation.onAnimation)
        
    };

    function handleOnclickDown(){
        if(count > 0) {
            setCount(prev => prev - 1);
        }
        console.log(Answers)
    };

    function handleScroll(e: React.WheelEvent){
        setScrollY( prev => prev + Math.ceil(e.deltaY * 0.0001))
    }

    return (
            <div style={{ maxHeight: '400px', overflowY: 'auto' }} onWheel={(e)=>handleScroll(e)}> 
                <div className={status} >
                    <h1 onClick={()=> console.log("Props", props)}>{scrollY}</h1>
                        {components[count]} 
                    <div className='flex flex-row' >
                        <h4 onClick={handleOnclickUp} className='m-2'>down</h4>
                        <h4 onClick={handleOnclickDown}>up</h4>
                    </div>
                </div>
                <button type='submit' onClick={() => handleSubmitForm(props.form_id, props.account_id, formRef)}>Submit</button>
            </div>
	)
    
};

