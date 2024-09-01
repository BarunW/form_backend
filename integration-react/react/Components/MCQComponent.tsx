import React, { } from "react";
import {ReactElement } from "react";
import { ComponentsParams } from "../Types";
import {  InsertAnswer } from "../Form";
function MakeChoice(question: any, qId : Number) : Array<ReactElement>{

    const options: Object = question.data.choices
    let choices : ReactElement[] = []

    Object.keys(options).forEach(key => { 
        choices.push(
            <div id="${key}" 
            className="p-2 border border-sky-500 rounded" 
            onClick={(event) => HandleMCQAnswer(event, qId)}>{options[key]["label"]}
            </div>
            )
    })
    return choices
}

function HandleMCQAnswer(e:any, qId: Number) { 
    InsertAnswer(qId, e.target.innerText)
}

export function MCQ(params : ComponentsParams){
    const { questionNumber, question, questionUUID} = params
    return(
    <div>
        <div className="flex">
        <span>{questionNumber}</span>
        <span>{" -> "}</span>
        <span>{question.data.question}</span>
        </div>
        <span>{question.data.description}</span>
        {MakeChoice(question, questionUUID)}
    </div>
    )
}
