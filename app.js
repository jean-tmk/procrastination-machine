const games=[
 {id:"planes",title:"Catch the runaway memo flock",note:"Move the inbox and catch six falling paper planes."},
 {id:"clock",title:"Stop the clock at precisely later",note:"Freeze the runaway hand inside the glowing slice three times."},
 {id:"snail",title:"Guide a snail through the inbox maze",note:"Find a safe route from the outbox to the tiny lamp."},
 {id:"pairs",title:"Match the secret office friendships",note:"Turn over the illustrated tiles and find all four pairs."},
 {id:"thread",title:"Untangle the red-thread conspiracy",note:"Trace the pins in order without losing the plot."},
 {id:"rhythm",title:"Play the desk-lamp rhythm",note:"Watch the office orchestra, then echo its six-beat pattern."},
 {id:"cabinet",title:"Sort the impossible filing cabinet",note:"Read each clue and file six wandering records correctly."}
];
const art=["clock","lamp","snail","paper plane","key","pen","paper stack","thread","chair","clip","envelope","file cabinet"];
const S={task:"",q:[],done:new Set(),start:0,active:0,tick:null,focus:null,cleanup:null};
const $=selector=>document.querySelector(selector);

setInterval(()=>$("#clock").textContent=new Date().toLocaleTimeString([],{hour:"2-digit",minute:"2-digit"}),1000);
document.querySelectorAll("nav button").forEach(button=>button.onclick=()=>$("#task").value=button.textContent);
function mix(items){return [...items].sort(()=>Math.random()-.5)}
function begin(task){
 S.task=task;S.done.clear();S.start=Date.now();S.q=mix(games).slice(0,5);
 $("#welcome").hidden=true;$("#play").hidden=false;$("#realTask").textContent=task;
 render();clearInterval(S.tick);S.tick=setInterval(metrics,1000);metrics();
}
document.querySelector("form").onsubmit=event=>{event.preventDefault();const task=$("#task").value.trim();if(task)begin(task)};

function render(){
 const quests=$("#quests");quests.innerHTML="";
 S.q.forEach((game,index)=>{
  const button=document.createElement("button"),locked=index&&!S.done.has(index-1),sprite=(games.findIndex(g=>g.id===game.id)+3)%art.length;
  button.className=`quest ${S.done.has(index)?"done":""} ${locked?"locked":""}`;
  button.innerHTML=`<small>DETOUR 0${index+1}</small><strong>${game.title}</strong><b class="s${sprite}" aria-label="${art[sprite]}">${S.done.has(index)?"✓":art[sprite]}</b><span>${S.done.has(index)?"ACTUALLY PERFORMED":"OPEN THE LITTLE DOOR"}</span>`;
  button.onclick=()=>!locked&&!S.done.has(index)&&openGame(index);
  quests.append(button);
 });
}

function openGame(index){
 if(S.cleanup)S.cleanup();
 S.active=index;const game=S.q[index],activity=$("#activity");
 $("#gameTitle").textContent=game.title;$("#instruction").textContent=game.note;
 activity.className=`mini-game ${game.id}`;activity.innerHTML="";
 const builders={planes:buildPlanes,clock:buildClock,snail:buildSnail,pairs:buildPairs,thread:buildThread,rhythm:buildRhythm,cabinet:buildCabinet};
 builders[game.id](activity);
 $("#game").showModal();
}

function buildPlanes(stage){
 stage.innerHTML='<div class="sky"></div><div class="catcher"><span>INBOX</span></div><p class="game-readout">CAUGHT <b>0</b> / 6</p>';
 const sky=stage.querySelector(".sky"),catcher=stage.querySelector(".catcher");let caught=0;
 for(let i=0;i<8;i++){
  const plane=document.createElement("button");plane.className="falling-plane sprite s3";plane.setAttribute("aria-label","falling paper plane");
  plane.style.setProperty("--x",`${5+Math.random()*82}%`);plane.style.setProperty("--delay",`${-Math.random()*5}s`);plane.style.setProperty("--speed",`${3.3+Math.random()*2}s`);sky.append(plane);
 }
 const move=event=>{const r=sky.getBoundingClientRect();catcher.style.left=`${Math.max(8,Math.min(92,(event.clientX-r.left)/r.width*100))}%`};
 sky.addEventListener("pointermove",move);
 const timer=setInterval(()=>{const cr=catcher.getBoundingClientRect();stage.querySelectorAll(".falling-plane").forEach(plane=>{const pr=plane.getBoundingClientRect();if(pr.bottom>cr.top&&pr.top<cr.bottom&&pr.right>cr.left&&pr.left<cr.right){caught++;stage.querySelector(".game-readout b").textContent=caught;plane.style.animationDelay="0s";plane.classList.add("caught");setTimeout(()=>plane.classList.remove("caught"),220);prog(caught,6)}})},80);
 S.cleanup=()=>clearInterval(timer);prog(0,6);
}

function buildClock(stage){
 stage.innerHTML='<div class="clock-face"><i class="target"></i><b class="clock-hand"></b><span>LATER</span></div><button class="clock-stop">FREEZE TIME</button><p class="game-readout">PERFECT STOPS <b>0</b> / 3</p>';
 const hand=stage.querySelector(".clock-hand"),button=stage.querySelector(".clock-stop");let hits=0,start=performance.now();
 button.onclick=()=>{const angle=((performance.now()-start)/2200*360)%360,distance=Math.min(Math.abs(angle),360-Math.abs(angle));hand.classList.add("paused");if(distance<38){hits++;stage.classList.add("time-hit");stage.querySelector(".game-readout b").textContent=hits;prog(hits,3)}else stage.classList.add("time-miss");setTimeout(()=>{stage.classList.remove("time-hit","time-miss");hand.classList.remove("paused");start=performance.now()},420)};
 prog(0,3);
}

function buildSnail(stage){
 const blocks=new Set([6,7,12,17]),start=20,goal=4;let position=start,moves=0;
 stage.innerHTML='<div class="maze-board"></div><div class="maze-controls"><button data-d="-5">UP</button><button data-d="-1">LEFT</button><button data-d="1">RIGHT</button><button data-d="5">DOWN</button></div><p class="game-readout">MOVES <b>0</b></p>';
 const board=stage.querySelector(".maze-board");
 for(let i=0;i<25;i++){const cell=document.createElement("div");cell.className=`maze-cell ${blocks.has(i)?"blocked":""} ${i===goal?"goal sprite s1":""}`;cell.dataset.cell=i;board.append(cell)}
 const draw=()=>{board.querySelectorAll(".snail-piece").forEach(x=>x.remove());const snail=document.createElement("i");snail.className="snail-piece sprite s2";board.children[position].append(snail)};
 stage.querySelectorAll(".maze-controls button").forEach(button=>button.onclick=()=>{const d=+button.dataset.d,next=position+d,sameRow=d===1||d===-1?Math.floor(next/5)===Math.floor(position/5):true;if(next>=0&&next<25&&sameRow&&!blocks.has(next)){position=next;moves++;stage.querySelector(".game-readout b").textContent=moves;draw();prog(position===goal?1:0,1)}});
 draw();prog(0,1);
}

function buildPairs(stage){
 const deck=mix([0,0,1,1,2,2,3,3]);let open=[],found=0,locked=false;
 stage.innerHTML='<div class="pair-grid"></div><p class="game-readout">FRIENDSHIPS <b>0</b> / 4</p>';
 const grid=stage.querySelector(".pair-grid");
 deck.forEach((sprite,index)=>{const card=document.createElement("button");card.className="pair-card";card.dataset.sprite=sprite;card.dataset.index=index;card.innerHTML=`<i class="sprite s${sprite}"></i><span>?</span>`;card.onclick=()=>{if(locked||card.classList.contains("open")||card.classList.contains("matched"))return;card.classList.add("open");open.push(card);if(open.length===2){locked=true;if(open[0].dataset.sprite===open[1].dataset.sprite){open.forEach(x=>x.classList.add("matched"));found++;stage.querySelector(".game-readout b").textContent=found;open=[];locked=false;prog(found,4)}else setTimeout(()=>{open.forEach(x=>x.classList.remove("open"));open=[];locked=false},650)}};grid.append(card)});
 prog(0,4);
}

function buildThread(stage){
 const points=[[11,72],[25,20],[42,66],[58,28],[73,76],[88,25]];let next=0;
 stage.innerHTML='<div class="thread-board"><svg viewBox="0 0 100 100" preserveAspectRatio="none"><polyline points=""></polyline></svg></div><p class="game-readout">CONNECTED <b>0</b> / 6</p>';
 const board=stage.querySelector(".thread-board"),line=board.querySelector("polyline"),connected=[];
 points.forEach((point,index)=>{const pin=document.createElement("button");pin.className=`thread-pin sprite s${[4,9,10,5,7,6][index]}`;pin.style.left=point[0]+"%";pin.style.top=point[1]+"%";pin.dataset.n=index;pin.setAttribute("aria-label",`pin ${index+1}`);const number=document.createElement("b");number.textContent=index+1;pin.append(number);pin.onclick=()=>{if(index!==next){board.classList.add("wrong-thread");setTimeout(()=>board.classList.remove("wrong-thread"),300);return}pin.classList.add("linked");connected.push(point.join(","));line.setAttribute("points",connected.join(" "));next++;stage.querySelector(".game-readout b").textContent=next;prog(next,6)};board.append(pin)});
 prog(0,6);
}

function buildRhythm(stage){
 const sequence=[0,2,3,1,2,0];let position=0,accepting=false;
 stage.innerHTML='<div class="rhythm-stage"><button class="sprite s1" aria-label="lamp"></button><button class="sprite s0" aria-label="clock"></button><button class="sprite s5" aria-label="pen"></button><button class="sprite s9" aria-label="clip"></button></div><button class="replay">SHOW ME THE RHYTHM</button><p class="game-readout">BEATS <b>0</b> / 6</p>';
 const pads=[...stage.querySelectorAll(".rhythm-stage button")];
 const show=()=>{accepting=false;position=0;stage.querySelector(".game-readout b").textContent=0;sequence.forEach((n,i)=>setTimeout(()=>{pads[n].classList.add("sing");setTimeout(()=>pads[n].classList.remove("sing"),280);if(i===sequence.length-1)setTimeout(()=>accepting=true,360)},i*430))};
 pads.forEach((pad,index)=>pad.onclick=()=>{if(!accepting)return;pad.classList.add("sing");setTimeout(()=>pad.classList.remove("sing"),180);if(index===sequence[position]){position++;stage.querySelector(".game-readout b").textContent=position;prog(position,6)}else{position=0;stage.querySelector(".game-readout b").textContent=0}});
 stage.querySelector(".replay").onclick=show;show();prog(0,6);
}

function buildCabinet(stage){
 const records=mix([{n:"MIDNIGHT IDEAS",slot:"A"},{n:"THINGS FOR LATER",slot:"B"},{n:"TINY EMERGENCIES",slot:"C"},{n:"DREAM RECEIPTS",slot:"A"},{n:"MISSING MINUTES",slot:"B"},{n:"POLITE PANICS",slot:"C"}]);let selected=null,filed=0;
 stage.innerHTML='<div class="records"></div><div class="drawers"><button data-slot="A">AFTER DARK</button><button data-slot="B">NOT YET</button><button data-slot="C">ODDLY URGENT</button></div><p class="game-readout">FILED <b>0</b> / 6</p>';
 const recordsEl=stage.querySelector(".records");
 records.forEach(record=>{const card=document.createElement("button");card.className="record-card";card.dataset.slot=record.slot;card.innerHTML=`<span>${record.n}</span><small>CLUE ${record.slot}</small>`;card.onclick=()=>{recordsEl.querySelectorAll("button").forEach(x=>x.classList.remove("selected"));card.classList.add("selected");selected=card};recordsEl.append(card)});
 stage.querySelectorAll(".drawers button").forEach(drawer=>drawer.onclick=()=>{if(!selected)return;if(selected.dataset.slot===drawer.dataset.slot){selected.classList.add("filed");selected=null;filed++;stage.querySelector(".game-readout b").textContent=filed;prog(filed,6)}else{selected.classList.add("wrong");setTimeout(()=>selected&&selected.classList.remove("wrong"),350)}});
 prog(0,6);
}

function prog(current,total){$("#progress").textContent=`${current} / ${total}`;$("#bar").style.width=Math.min(100,current/total*100)+"%";if(current>=total)setTimeout(complete,500)}
function complete(){if(S.cleanup){S.cleanup();S.cleanup=null}S.done.add(S.active);$("#game").close();render();metrics();$("#toast").classList.add("show");setTimeout(()=>$("#toast").classList.remove("show"),1700)}
function metrics(){const seconds=Math.floor((Date.now()-S.start)/1000);$("#lost").textContent=`${String(seconds/60|0).padStart(2,"0")}:${String(seconds%60).padStart(2,"0")}`;$("#count").textContent=S.done.size;$("#fake").textContent=Math.min(99,S.done.size*18+(seconds/15|0))+"%"}
$(".x").onclick=()=>{if(S.cleanup){S.cleanup();S.cleanup=null}$("#game").close()};
$("#shuffle").onclick=()=>begin(S.task);$("#restart").onclick=()=>{clearInterval(S.tick);$("#play").hidden=true;$("#welcome").hidden=false};
$("#actually").onclick=()=>{$("#focusTask").textContent=S.task;$("#focus").showModal()};$(".fx").onclick=()=>$("#focus").close();
$("#five").onclick=()=>{clearInterval(S.focus);let seconds=300;S.focus=setInterval(()=>{seconds--;$("#timer").textContent=`${String(seconds/60|0).padStart(2,"0")}:${String(seconds%60).padStart(2,"0")}`;if(!seconds)clearInterval(S.focus)},1000)};
$("#finished").onclick=()=>{$("#focus").close();$("#restart").click();$("#toast").textContent="WAIT… YOU ACTUALLY DID IT?!";$("#toast").classList.add("show")};
