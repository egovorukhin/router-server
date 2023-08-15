import{_ as w,u as k,f as i,g as V,c as C,r as l,o as S,a as I,d as s,b as e,F as T,h as B,p as L,i as j,e as U}from"./index-254ebb3c.js";import{C as N}from"./CrudService-9e941e24.js";const a=n=>(L("data-v-f680840d"),n=n(),j(),n),A={class:"surface-ground flex align-items-center justify-content-center min-h-screen min-w-screen overflow-hidden"},F={class:"flex flex-column align-items-center justify-content-center"},P=["src"],$={style:{"border-radius":"56px",padding:"0.3rem",background:"linear-gradient(180deg, var(--primary-color) 10%, rgba(33, 150, 243, 0) 30%)"}},E={class:"w-full surface-card py-8 px-5 sm:px-8",style:{"border-radius":"53px"}},J=a(()=>e("div",{class:"text-center mb-5"},[e("div",{class:"text-900 text-3xl font-medium mb-3"},"ALAcall Contact Center"),e("span",{class:"text-600 font-medium"},"Войдите для продолжения")],-1)),M=a(()=>e("label",{for:"email1",class:"block text-900 text-xl font-medium mb-2"},"Имя",-1)),O=a(()=>e("label",{for:"password1",class:"block text-900 font-medium text-xl mb-2"},"Пароль",-1)),q={class:"flex align-items-center justify-content-between mb-5 gap-5"},z={class:"flex align-items-center"},D=a(()=>e("label",{for:"rememberme1"},"Запомнить меня",-1)),G=a(()=>e("a",{class:"font-medium no-underline ml-2 text-right cursor-pointer",style:{color:"var(--primary-color)"}},"Забыли пароль?",-1)),H={__name:"Login",setup(n){const{layoutConfig:u}=k(),c=i(""),d=i(""),m=i(!1),p=V(),_=new N(""),f=C(()=>`layout/images/${u.darkTheme.value?"logo-white":"logo-dark"}.svg`),g=()=>{_.auth(c.value,d.value).then(r=>{r.status===200?r.json().then(t=>{t&&(localStorage.setItem("session",JSON.stringify(t)),U.push("/"))}):r.text().then(t=>{p.add({severity:"error",summary:"Ошибка",detail:t,life:3e3})})})};return(r,t)=>{const x=l("Toast"),v=l("InputText"),b=l("Password"),h=l("Checkbox"),y=l("Button");return S(),I(T,null,[s(x),e("div",A,[e("div",F,[e("img",{src:f.value,alt:"Sakai logo",class:"mb-5 w-6rem flex-shrink-0"},null,8,P),e("div",$,[e("div",E,[J,e("div",null,[M,s(v,{id:"email1",type:"text",placeholder:"Учётная запись",class:"w-full md:w-30rem mb-5",style:{padding:"1rem"},modelValue:c.value,"onUpdate:modelValue":t[0]||(t[0]=o=>c.value=o)},null,8,["modelValue"]),O,s(b,{id:"password1",feedback:!1,modelValue:d.value,"onUpdate:modelValue":t[1]||(t[1]=o=>d.value=o),placeholder:"Пароль",toggleMask:!0,class:"w-full mb-3",inputClass:"w-full",inputStyle:{padding:"1rem"}},null,8,["modelValue"]),e("div",q,[e("div",z,[s(h,{modelValue:m.value,"onUpdate:modelValue":t[2]||(t[2]=o=>m.value=o),id:"rememberme1",binary:"",class:"mr-2"},null,8,["modelValue"]),D]),G]),s(y,{label:"Войти",class:"w-full p-3 text-xl",onClick:g})])])])])]),s(B,{simple:""})],64)}}},R=w(H,[["__scopeId","data-v-f680840d"]]);export{R as default};
