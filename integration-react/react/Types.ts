export interface PropsType {
    account_id      : string,
    form_id         : string;
    questions       : any;
    form_settings   : Object;
    answers         : Object;
    
}

export interface ComponentsParams {
    questionUUID    : Number;
    questionNumber  : Number;
    question        : any;
    formSettings    : Object;
}


export interface inputProperties{
    label       : string;
    name        : string;
    placeholder : string;
    type        : string;
}




