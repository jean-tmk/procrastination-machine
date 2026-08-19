const games=[
 {id:"planes",icon:3,title:"Catch the runaway memo flock",note:"Move the inbox and catch six falling paper planes."},
 {id:"clock",icon:0,title:"Stop the clock at precisely later",note:"Freeze the orbiting red marker inside the golden window three times."},
 {id:"snail",icon:2,title:"Guide a snail through the inbox maze",note:"Explore the shifting maze and guide the snail to the red thread."},
 {id:"pairs",icon:9,title:"Match the secret office friendships",note:"Turn over the illustrated tiles and find all four pairs."},
 {id:"thread",icon:7,title:"Untangle the red-thread conspiracy",note:"Trace the pins in order without losing the plot."},
 {id:"rhythm",icon:1,title:"Play the desk-lamp rhythm",note:"Watch the office orchestra, then echo its six-beat pattern."},
 {id:"cabinet",icon:11,title:"Sort the impossible filing cabinet",note:"Read each clue and file six wandering records correctly."},
 {id:"typewriter",icon:5,title:"Repair the forgetful typewriter",note:"Catch the missing letters and rebuild three mysteriously damaged words."},
 {id:"spotlight",icon:1,title:"Search the desk after midnight",note:"Sweep the lamp across the dark desk and uncover four hidden objects."},
 {id:"balance",icon:6,title:"Balance the tower of unfinished business",note:"Stack seven drifting papers without letting the tower tip."},
 {id:"switchboard",icon:4,title:"Reroute the office daydream",note:"Rotate the connections until every glowing thought reaches the red phone."},
 {id:"stamp",icon:9,title:"Stamp the forms before they escape",note:"Catch and certify six glowing forms before the desk changes its mind."},
 {id:"shredder",icon:6,title:"Feed the shredder only bad ideas",note:"Destroy the five terrible ideas while protecting the suspiciously useful ones."},
 {id:"cursor",icon:5,title:"Trap the runaway cursor",note:"Catch the blinking cursor eight times before it files itself elsewhere."},
 {id:"shelves",icon:11,title:"Arrange the cabinet of tiny excuses",note:"Select the five excuses from the shortest delay to the most enormous."}
];
const art=["clock","lamp","snail","paper plane","key","pen","paper stack","thread","chair","clip","envelope","file cabinet"];
const S={task:"",q:[],remaining:[],done:new Set(),start:0,active:0,tick:null,focus:null,cleanup:null};
const $=selector=>document.querySelector(selector);

setInterval(()=>$("#clock").textContent=new Date().toLocaleTimeString([],{hour:"2-digit",minute:"2-digit"}),1000);
document.querySelectorAll("nav button").forEach(button=>button.onclick=()=>$("#task").value=button.textContent);
function mix(items){return [...items].sort(()=>Math.random()-.5)}
function loadCycle(){try{const saved=JSON.parse(localStorage.getItem("pm-game-cycle-v2")||"[]"),valid=new Set(games.map(game=>game.id));return saved.filter((id,index)=>valid.has(id)&&saved.indexOf(id)===index)}catch{return[]}}
function refillCycle(){const previous=new Set(S.q.map(game=>game.id)),fresh=mix(games.map(game=>game.id)),notRecent=fresh.filter(id=>!previous.has(id)),recent=fresh.filter(id=>previous.has(id));S.remaining=[...notRecent,...recent]}
function dealRound(){if(!S.remaining.length)S.remaining=loadCycle();if(S.remaining.length<5)refillCycle();const ids=S.remaining.splice(0,5);localStorage.setItem("pm-game-cycle-v2",JSON.stringify(S.remaining));S.q=ids.map(id=>games.find(game=>game.id===id));S.done.clear()}
function begin(task){
 S.task=task;S.start=Date.now();dealRound();
 $("#welcome").hidden=true;$("#play").hidden=false;$("#realTask").textContent=task;
 render();clearInterval(S.tick);S.tick=setInterval(metrics,1000);metrics();
}
document.querySelector("form").onsubmit=event=>{event.preventDefault();const task=$("#task").value.trim();if(task)begin(task)};

function render(){
 const quests=$("#quests");quests.innerHTML="";
 S.q.forEach((game,index)=>{
  const button=document.createElement("button"),locked=index&&!S.done.has(index-1),sprite=game.icon;
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
 const builders={planes:buildPlanes,clock:buildClock,snail:buildSnail,pairs:buildPairs,thread:buildThread,rhythm:buildRhythm,cabinet:buildCabinet,typewriter:buildTypewriter,spotlight:buildSpotlight,balance:buildBalance,switchboard:buildSwitchboard,stamp:buildStamp,shredder:buildShredder,cursor:buildCursor,shelves:buildShelves};
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
 stage.innerHTML='<div class="clock-face sprite s0"><div class="clock-center"><i class="target"></i><b class="clock-hand"></b></div><span>PRECISELY LATER</span></div><button class="clock-stop">FREEZE THE ORBIT</button><p class="game-readout">PERFECT STOPS <b>0</b> / 3</p>';
 const hand=stage.querySelector(".clock-hand"),button=stage.querySelector(".clock-stop");let hits=0,start=performance.now();
 button.onclick=()=>{const angle=((performance.now()-start)/2200*360)%360,distance=Math.min(Math.abs(angle),360-Math.abs(angle));hand.classList.add("paused");if(distance<38){hits++;stage.classList.add("time-hit");stage.querySelector(".game-readout b").textContent=hits;prog(hits,3)}else stage.classList.add("time-miss");setTimeout(()=>{stage.classList.remove("time-hit","time-miss");hand.classList.remove("paused");start=performance.now()},420)};
 prog(0,3);
}

function buildSnail(stage){
 const size=15,cells=Array(size*size).fill("#"),start=size+1,stack=[start],directions=[-size,size,-1,1];cells[start]=".";
 while(stack.length){const current=stack[stack.length-1],x=current%size,y=Math.floor(current/size),options=mix(directions).filter(d=>{const next=current+d*2,nx=next%size,ny=Math.floor(next/size);return nx>0&&nx<size-1&&ny>0&&ny<size-1&&cells[next]==="#"&&((d===1||d===-1)?Math.abs(nx-x)===2:Math.abs(ny-y)===2)});if(!options.length){stack.pop();continue}const d=options[0],next=current+d*2;cells[current+d]=".";cells[next]=".";stack.push(next)}
 const distance=Array(size*size).fill(-1),queue=[start];distance[start]=0;for(let head=0;head<queue.length;head++){const current=queue[head];directions.forEach(d=>{const next=current+d,sameRow=d===1||d===-1?Math.floor(next/size)===Math.floor(current/size):true;if(next>=0&&next<cells.length&&sameRow&&cells[next]==="."&&distance[next]<0){distance[next]=distance[current]+1;queue.push(next)}})}const goal=distance.indexOf(Math.max(...distance)),blocks=new Set(cells.map((cell,index)=>cell==="#"?index:-1).filter(index=>index>=0));let position=start,moves=0;const seen=new Set();
 stage.innerHTML=`<div class="maze-board" style="--maze-size:${size}"></div><div class="maze-controls"><button class="up" data-d="-${size}" aria-label="move up">UP</button><button class="left" data-d="-1" aria-label="move left">LEFT</button><i></i><button class="right" data-d="1" aria-label="move right">RIGHT</button><button class="down" data-d="${size}" aria-label="move down">DOWN</button></div><p class="game-readout">MOVES <b>0</b> · THREAD DISTANCE <strong>${distance[goal]}</strong></p>`;
 const board=stage.querySelector(".maze-board");
 for(let i=0;i<cells.length;i++){const cell=document.createElement("div");cell.className=`maze-cell ${blocks.has(i)?"blocked":""} ${i===goal?"goal sprite s7":""}`;board.append(cell)}
 const draw=()=>{board.querySelectorAll(".snail-piece").forEach(x=>x.remove());const px=position%size,py=Math.floor(position/size);for(let y=Math.max(0,py-2);y<=Math.min(size-1,py+2);y++)for(let x=Math.max(0,px-2);x<=Math.min(size-1,px+2);x++)if(Math.abs(px-x)+Math.abs(py-y)<=3)seen.add(y*size+x);[...board.children].forEach((cell,index)=>cell.classList.toggle("seen",seen.has(index)));const snail=document.createElement("i");snail.className="snail-piece sprite s2";board.children[position].append(snail)};
 stage.querySelectorAll(".maze-controls button").forEach(button=>button.onclick=()=>{const d=+button.dataset.d,next=position+d,sameRow=d===1||d===-1?Math.floor(next/size)===Math.floor(position/size):true;if(next>=0&&next<cells.length&&sameRow&&!blocks.has(next)){position=next;moves++;stage.querySelector(".game-readout b").textContent=moves;draw();prog(position===goal?1:0,1)}});
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
 const points=[[13,22],[31,78],[46,28],[61,72],[76,18],[88,64]],icons=[7,4,10,5,9,6],route=[1,4,0,2,5,3];let next=0;
 stage.innerHTML='<div class="thread-clue">KEY <span>then</span> CLIP <span>then</span> THREAD <span>then</span> ENVELOPE <span>then</span> PAPER <span>then</span> PEN</div><div class="thread-board"><svg viewBox="0 0 100 100" preserveAspectRatio="none"><polyline points=""></polyline></svg></div><p class="game-readout">CONNECTED <b>0</b> / 6</p>';
 const board=stage.querySelector(".thread-board"),line=board.querySelector("polyline"),connected=[];
 points.forEach((point,index)=>{const pin=document.createElement("button");pin.className=`thread-pin sprite s${icons[index]}`;pin.style.left=point[0]+"%";pin.style.top=point[1]+"%";pin.setAttribute("aria-label",art[icons[index]]);pin.onclick=()=>{if(index!==route[next]){board.classList.add("wrong-thread");setTimeout(()=>board.classList.remove("wrong-thread"),300);return}pin.classList.add("linked");connected.push(point.join(","));line.setAttribute("points",connected.join(" "));next++;stage.querySelector(".game-readout b").textContent=next;prog(next,6)};board.append(pin)});
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

function buildTypewriter(stage){
 const puzzles=[{word:"LATER",blank:1,letter:"A"},{word:"MAYBE",blank:1,letter:"A"},{word:"PAUSE",blank:2,letter:"U"}];let round=0;
 stage.innerHTML='<div class="typewriter-machine"><div class="paper-word"></div><div class="type-keys"></div></div><p class="type-note">The typewriter misplaced one letter.</p><p class="game-readout">WORDS REPAIRED <b>0</b> / 3</p>';
 const word=stage.querySelector(".paper-word"),keys=stage.querySelector(".type-keys"),note=stage.querySelector(".type-note");
 const show=()=>{const puzzle=puzzles[round];word.textContent=[...puzzle.word].map((letter,index)=>index===puzzle.blank?"_":letter).join(" ");keys.innerHTML="";mix(["A","E","I","O","U"]).forEach(letter=>{const key=document.createElement("button");key.textContent=letter;key.onclick=()=>{if(letter===puzzle.letter){word.textContent=[...puzzle.word].join(" ");note.textContent="The sentence can breathe again.";round++;stage.querySelector(".game-readout b").textContent=round;prog(round,3);if(round<3)setTimeout(show,550)}else{key.classList.add("wrong");note.textContent="That letter belongs to a different excuse.";setTimeout(()=>key.classList.remove("wrong"),300)}};keys.append(key)})};show();prog(0,3);
}
function buildSpotlight(stage){
 const positions=[[14,18],[72,16],[25,70],[76,68]];let found=0;
 stage.innerHTML='<div class="spotlight-desk"></div><p class="game-readout">DISCOVERED <b>0</b> / 4</p>';
 const desk=stage.querySelector(".spotlight-desk");
 positions.forEach((position,index)=>{const object=document.createElement("button");object.className=`hidden-object sprite s${[4,9,10,2][index]}`;object.style.left=position[0]+"%";object.style.top=position[1]+"%";object.setAttribute("aria-label",art[[4,9,10,2][index]]);object.onclick=()=>{if(object.classList.contains("found"))return;object.classList.add("found");found++;stage.querySelector(".game-readout b").textContent=found;prog(found,4)};desk.append(object)});
 desk.onpointermove=event=>{const rect=desk.getBoundingClientRect();desk.style.setProperty("--mx",event.clientX-rect.left+"px");desk.style.setProperty("--my",event.clientY-rect.top+"px")};prog(0,4);
}
function buildBalance(stage){
 let level=0,last=50,start=performance.now();
 stage.innerHTML='<div class="balance-scene"><div class="paper-stack"></div><div class="drifting-paper">NEXT FILE</div></div><button class="drop-paper">DROP THE PAPER</button><p class="balance-note">Keep each page close to the center of the stack.</p><p class="game-readout">BALANCED <b>0</b> / 7</p>';
 const scene=stage.querySelector(".balance-scene"),stack=stage.querySelector(".paper-stack"),moving=stage.querySelector(".drifting-paper"),note=stage.querySelector(".balance-note");
 const timer=setInterval(()=>{const x=50+Math.sin((performance.now()-start)/520)*39;moving.style.left=x+"%";moving.dataset.x=x},20);S.cleanup=()=>clearInterval(timer);
 stage.querySelector(".drop-paper").onclick=()=>{const x=+moving.dataset.x||50;if(level&&Math.abs(x-last)>24){stack.innerHTML="";level=0;last=50;scene.classList.add("topple");note.textContent="The unfinished business toppled. Begin the pile again.";stage.querySelector(".game-readout b").textContent=0;prog(0,7);setTimeout(()=>scene.classList.remove("topple"),500);return}const page=document.createElement("i");page.style.left=x+"%";page.style.bottom=14+level*20+"px";stack.append(page);last=x;level++;stage.querySelector(".game-readout b").textContent=level;note.textContent=level<7?"Suspiciously stable. Add another.":"A monument to doing it later.";prog(level,7)};prog(0,7);
}
function buildSwitchboard(stage){
 const rotations=mix([1,2,3,1,2,3,1,2,3]);let solved=0;
 stage.innerHTML='<div class="switch-grid"></div><p class="switch-note">Click each brass plate until its red line runs straight upward.</p><p class="game-readout">ALIGNED <b>0</b> / 9</p>';
 const grid=stage.querySelector(".switch-grid");
 rotations.forEach((rotation,index)=>{const tile=document.createElement("button");tile.className="switch-tile";tile.style.setProperty("--turn",rotation*90+"deg");tile.innerHTML=`<i></i><span>${String(index+1).padStart(2,"0")}</span>`;tile.onclick=()=>{rotations[index]=(rotations[index]+1)%4;tile.style.setProperty("--turn",rotations[index]*90+"deg");solved=rotations.filter(value=>value===0).length;stage.querySelector(".game-readout b").textContent=solved;tile.classList.toggle("aligned",rotations[index]===0);prog(solved,9)};grid.append(tile)});prog(0,9);
}

function buildStamp(stage){
 let active=Math.floor(Math.random()*9),stamped=0;
 stage.innerHTML='<div class="stamp-desk"></div><p class="stamp-note">The glowing form is currently willing to be approved.</p><p class="game-readout">CERTIFIED <b>0</b> / 6</p>';
 const desk=stage.querySelector(".stamp-desk"),draw=()=>{[...desk.children].forEach((form,index)=>form.classList.toggle("active",index===active))};
 for(let index=0;index<9;index++){const form=document.createElement("button");form.innerHTML=`<span>FORM ${String(index+1).padStart(2,"0")}</span><i></i>`;form.onclick=()=>{if(index!==active){form.classList.add("rejected");stage.querySelector(".stamp-note").textContent="That form was not emotionally ready.";setTimeout(()=>form.classList.remove("rejected"),300);return}form.classList.add("stamped");stamped++;stage.querySelector(".game-readout b").textContent=stamped;stage.querySelector(".stamp-note").textContent="CERTIFIED UNNECESSARY.";prog(stamped,6);do active=Math.floor(Math.random()*9);while(active===index);setTimeout(()=>{form.classList.remove("stamped");draw()},220)};desk.append(form)}draw();prog(0,6);
}

function buildShredder(stage){
 const files=mix([{t:"REPLY ALL FOREVER",bad:1},{t:"MEETING ABOUT MEETINGS",bad:1},{t:"PRINT THE INTERNET",bad:1},{t:"URGENT FONT EMERGENCY",bad:1},{t:"RENAME EVERY FOLDER",bad:1},{t:"SAVE THE ACTUAL WORK",bad:0},{t:"CALL SOMEONE KIND",bad:0},{t:"DRINK A GLASS OF WATER",bad:0}]);let shredded=0;
 stage.innerHTML='<div class="shredder-files"></div><div class="shredder-mouth"><span>BAD IDEAS ONLY</span><i></i></div><p class="shred-note">Choose carefully. The shredder has no undo button.</p><p class="game-readout">BAD IDEAS SHREDDED <b>0</b> / 5</p>';
 const tray=stage.querySelector(".shredder-files");files.forEach(file=>{const card=document.createElement("button");card.textContent=file.t;card.onclick=()=>{if(card.classList.contains("gone"))return;if(file.bad){card.classList.add("gone");shredded++;stage.querySelector(".game-readout b").textContent=shredded;stage.querySelector(".shred-note").textContent="An excellent administrative decision.";prog(shredded,5)}else{card.classList.add("protected");stage.querySelector(".shred-note").textContent="Protected! That one might accidentally help.";setTimeout(()=>card.classList.remove("protected"),450)}};tray.append(card)});prog(0,5);
}

function buildCursor(stage){
 let caught=0;
 stage.innerHTML='<div class="cursor-field"><button class="runaway-cursor" aria-label="catch the runaway cursor"><i></i></button><span>THIS DOCUMENT HAS 48 UNSAVED FEELINGS</span></div><p class="cursor-note">It becomes less cooperative every time.</p><p class="game-readout">CAUGHT <b>0</b> / 8</p>';
 const field=stage.querySelector(".cursor-field"),cursor=stage.querySelector(".runaway-cursor"),move=()=>{cursor.style.left=8+Math.random()*82+"%";cursor.style.top=12+Math.random()*72+"%";cursor.style.setProperty("--tilt",-25+Math.random()*50+"deg")};cursor.onclick=()=>{caught++;stage.querySelector(".game-readout b").textContent=caught;field.classList.add("cursor-caught");setTimeout(()=>field.classList.remove("cursor-caught"),160);cursor.style.scale=String(Math.max(.58,1-caught*.05));move();prog(caught,8)};move();prog(0,8);
}

function buildShelves(stage){
 const excuses=[{t:"one blink",h:112},{t:"tiny pause",h:136},{t:"short detour",h:160},{t:"extended wander",h:188},{t:"entire afternoon",h:218}],shuffled=mix(excuses);let next=0;
 stage.innerHTML='<div class="excuse-shelf"></div><p class="shelf-note">Begin with the smallest excuse and work upward.</p><p class="game-readout">SHELVED <b>0</b> / 5</p>';
 const shelf=stage.querySelector(".excuse-shelf");shuffled.forEach(excuse=>{const item=document.createElement("button");item.style.setProperty("--height",excuse.h+"px");item.innerHTML=`<i></i><span>${excuse.t}</span>`;item.onclick=()=>{if(item.classList.contains("placed"))return;if(excuse===excuses[next]){item.classList.add("placed");next++;stage.querySelector(".game-readout b").textContent=next;stage.querySelector(".shelf-note").textContent=next<5?"Correctly unreasonable. Keep going.":"A perfectly escalating wall of avoidance.";prog(next,5)}else{item.classList.add("wrong");stage.querySelector(".shelf-note").textContent="That excuse is wildly out of proportion.";setTimeout(()=>item.classList.remove("wrong"),350)}};shelf.append(item)});prog(0,5);
}

function prog(current,total){$("#progress").textContent=`${current} / ${total}`;$("#bar").style.width=Math.min(100,current/total*100)+"%";if(current>=total)setTimeout(complete,500)}
function complete(){if(S.cleanup){S.cleanup();S.cleanup=null}S.done.add(S.active);$("#game").close();render();metrics();$("#toast").classList.add("show");setTimeout(()=>$("#toast").classList.remove("show"),1700)}
function metrics(){const seconds=Math.floor((Date.now()-S.start)/1000);$("#lost").textContent=`${String(seconds/60|0).padStart(2,"0")}:${String(seconds%60).padStart(2,"0")}`;$("#count").textContent=S.done.size;$("#fake").textContent=Math.min(99,S.done.size*18+(seconds/15|0))+"%"}
$(".x").onclick=()=>{if(S.cleanup){S.cleanup();S.cleanup=null}$("#game").close()};
$("#shuffle").onclick=()=>begin(S.task);$("#restart").onclick=()=>{clearInterval(S.tick);$("#play").hidden=true;$("#welcome").hidden=false};
$("#actually").onclick=()=>{$("#focusTask").textContent=S.task;$("#focus").showModal()};$(".fx").onclick=()=>$("#focus").close();
$("#five").onclick=()=>{clearInterval(S.focus);let seconds=300;S.focus=setInterval(()=>{seconds--;$("#timer").textContent=`${String(seconds/60|0).padStart(2,"0")}:${String(seconds%60).padStart(2,"0")}`;if(!seconds)clearInterval(S.focus)},1000)};
$("#finished").onclick=()=>{$("#focus").close();$("#restart").click();$("#toast").textContent="WAIT… YOU ACTUALLY DID IT?!";$("#toast").classList.add("show")};
