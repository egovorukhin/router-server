import{_ as k,u as V,r as u,a as S,c as T,b as l,o as C,d as I,e as n,f as e,w as _,F as L,p as j,g as B,h as f}from"./index-5e528b76.js";const s=a=>(j("data-v-618c1367"),a=a(),B(),a),U={class:"surface-ground flex align-items-center justify-content-center min-h-screen min-w-screen overflow-hidden"},P={class:"flex flex-column align-items-center justify-content-center"},K={style:{"border-radius":"56px",padding:"0.3rem",background:"linear-gradient(180deg, var(--primary-color) 10%, rgba(33, 150, 243, 0) 30%)"}},N={class:"w-full surface-card py-8 px-5 sm:px-8",style:{"border-radius":"53px"}},$={class:"text-center mb-5"},A=["src"],F=s(()=>e("div",{class:"text-900 text-3xl font-medium mb-3"},[e("span",{id:"logo-alacall"},"ALAcall"),e("span",{style:{"font-weight":"normal"}},"SipProxy")],-1)),O=s(()=>e("span",{class:"text-600 font-medium"},"Войдите для продолжения",-1)),E=s(()=>e("label",{for:"email1",class:"block text-900 text-xl font-medium mb-2"},"Имя",-1)),J=s(()=>e("label",{for:"password1",class:"block text-900 font-medium text-xl mb-2"},"Пароль",-1)),M={class:"flex align-items-center justify-content-between mb-5 gap-5"},q={class:"flex align-items-center"},z=s(()=>e("label",{for:"rememberme1"},"Запомнить меня",-1)),D=s(()=>e("a",{class:"font-medium no-underline ml-2 text-right cursor-pointer",style:{color:"var(--primary-color)"}},"Забыли пароль?",-1)),G={__name:"Login",setup(a){const{layoutConfig:g}=V(),i=u(""),d=u(""),m=u(!1),p=S(),x=T(()=>`layout/images/${g.darkTheme.value?"logo-white":"logo-dark"}.svg`),c=()=>{fetch((window.apiUrl??"")+"/api/auth",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({username:i.value,password:d.value})}).then(r=>{if(r.ok){r.json().then(t=>{if(t){f.session.setSession(t),f.r.push("/");return}p.add({severity:"error",summary:"Ошибка",detail:"Невалидные данные",life:3e3})});return}r.text().then(t=>{console.error(t),p.add({severity:"error",summary:"Ошибка",detail:t,life:3e3})})})};return(r,t)=>{const v=l("Toast"),h=l("InputText"),y=l("Password"),b=l("Checkbox"),w=l("Button");return C(),I(L,null,[n(v),e("div",U,[e("div",P,[e("div",K,[e("div",N,[e("div",$,[e("img",{src:x.value,alt:"Image",height:"70",class:"mb-3"},null,8,A),F,O]),e("div",null,[E,n(h,{id:"email1",type:"text",placeholder:"Учётная запись",class:"w-full md:w-30rem mb-5",style:{padding:"1rem"},modelValue:i.value,"onUpdate:modelValue":t[0]||(t[0]=o=>i.value=o),onKeyup:t[1]||(t[1]=_(o=>c(),["enter"]))},null,8,["modelValue"]),J,n(y,{id:"password1",feedback:!1,modelValue:d.value,"onUpdate:modelValue":t[2]||(t[2]=o=>d.value=o),placeholder:"Пароль",onKeyup:t[3]||(t[3]=_(o=>c(),["enter"])),toggleMask:!0,class:"w-full mb-3",inputClass:"w-full",inputStyle:{padding:"1rem"}},null,8,["modelValue"]),e("div",M,[e("div",q,[n(b,{modelValue:m.value,"onUpdate:modelValue":t[4]||(t[4]=o=>m.value=o),id:"rememberme1",binary:"",class:"mr-2"},null,8,["modelValue"]),z]),D]),n(w,{id:"btnLogin",label:"Войти",class:"w-full p-3 text-xl",onClick:c})])])])])])],64)}}},Q=k(G,[["__scopeId","data-v-618c1367"]]);export{Q as default};
