"use strict";

/** @type {WebSocket} */
let ws;

/** @type {HTMLFormElement} */
let socketsForm;

/** @type {HTMLInputElement} */
let socketsFormDir;

/** @type {HTMLTextAreaElement} */
let socketsFormInput;

// /** @type {HTMLInputElement} */
// let socketsFormEnc;

/** @type {HTMLButtonElement} */
let socketsFormSubmit;

/** @type {HTMLDialogElement} */
let packetDialog;

/** @type {HTMLTableElement} */
let socketsTable;

window.addEventListener("load", () => {
	socketsForm = /** @type {HTMLFormElement} */ (getSockElemById("sockets-form"));
	socketsTable = /** @type {HTMLTableElement} */ (getSockElemById("sockets-table"));

	/** @type {string} */
	let wsScheme;
	if (document.location.protocol == "http:")
		wsScheme = "ws://";
	else
		wsScheme = "wss://";
	const wsURL = `${wsScheme}${window.location.hostname}:${window.location.port}/api/sockets`;
	ws = new WebSocket(wsURL);

	ws.addEventListener("error", () => {
		const errMsg = `Kļūme ar WebSocket (${wsURL}) izveidošanu`;
		showError(errMsg);
		throw Error(errMsg);
	});

	ws.addEventListener("open", () => {
		socketsForm.style.display = "block";
		socketsTable.style.visibility = "visible";
		socketsFormSubmit = /** @type {HTMLButtonElement} */ (getSockElemById("sockets-form-submit"));
		socketsFormSubmit.addEventListener("click", socketsHandleSubmit);
	});

	ws.addEventListener("message", (msg) => {
		let row = JSON.parse(msg.data);
		console.log("Recieved: ", row);

		if (row.Error)
			showError(row.ErrMsg);
		else
			addRowToTable(row);
	})
});

function addRowToTable(recvRow) {
	let tableRow = socketsTable.getElementsByTagName("tbody")[0].insertRow(0);
	addCellToRow(tableRow, recvRow.Sender.Name);
	addCellToRow(tableRow, recvRow.Sender.Addr);
	addCellToRow(tableRow, recvRow.Sender.Port);
	addDialogToRow(tableRow, recvRow);
	addCellToRow(tableRow, recvRow.Receiver.Name);
	addCellToRow(tableRow, recvRow.Receiver.Addr);
	addCellToRow(tableRow, recvRow.Receiver.Port);
}

function addCellToRow(tableRow, recvRowData) {
	tableRow.insertCell().appendChild(document.createTextNode(recvRowData));
}

function addDialogToRow(tableRow, recvRowData) {
	let link = document.createElement("a");
	link.innerText = "Skatīt";
	link.style.cursor = "pointer";
	tableRow.insertCell().appendChild(link);
	link.addEventListener("click", () => { showPacketDialog(recvRowData); });
}

function showPacketDialog(recvRowData) {
	packetDialog = /** @type {HTMLDialogElement} */ (getSockElemById("packet-dialog"));
	packetDialog.getElementsByTagName("p")[0].innerText = recvRowData.Wire.PacketDump;
	packetDialog.getElementsByTagName("button")[0].addEventListener("click", () => { packetDialog.close(); });
	packetDialog.showModal();
}

/** @returns {void} */
function Row(senderName, recvName, encrypted) {
	this.Error = false;
	this.ErrMsg = "";

	this.Sender = {
		Name: senderName,
		Addr: "",
		Port: 0,
	};

	this.Receiver = {
		Name: recvName,
		Addr: "",
		Port: 0,
	};

	this.Payload = {
		Encrypted: false,
		Msg: ""
	}

	if (encrypted)
		this.Payload.Encrypted = true;
}

/** @returns {void} */
function socketsHandleSubmit() {
	hideError();
	socketsFormInput = /** @type {HTMLTextAreaElement} */ (getSockElemById("sockets-form-input"));
	// socketsFormEnc = /** @type {HTMLInputElement} */ (getSockElemById("sockets-form-enc"));
	socketsFormDir = /** @type {HTMLInputElement} */ (document.querySelector("input[name='sockets-form-dir']:checked"));
	let row = new Row();

	if (socketsFormDir.value == "A2B") {
		row.Sender.Name = "A";
		row.Receiver.Name = "B";
	} else if (socketsFormDir.value == "B2A") {
		row.Sender.Name = "B";
		row.Receiver.Name = "A";
	} else {
		showError("Kļūme ar virziena izvēli");
	}
	// row.Payload.Encrypted = socketsFormEnc.checked;
	row.Payload.Msg = socketsFormInput.value;

	console.log("Sent: ", row);
	ws.send(JSON.stringify(row));
}

/**
 * @param {string} id 
 * @returns {HTMLElement}
 */
function getSockElemById(id) {
	let elem = document.getElementById(id);
	if (!elem) {
		const errMsg = `Kļūme, nevarēja atrast elementu ar ID "${id}"`;
		showError(errMsg);
		throw Error(errMsg);
	}
	return elem;
}

/**
 * @param {string} msg 
 * @returns {void}
 */
function showError(msg) {
	let socketErr = /** @type {HTMLDivElement} */ (document.getElementById("socketErr"));
	if (socketErr) {
		socketErr.remove();
	}

	let socketErrMsg = document.createElement("p");
	socketErrMsg.classList.add("errMsg");
	socketErrMsg.innerText = msg;

	socketErr = document.createElement("div");
	socketErr.classList.add("cw-center");
	socketErr.id = "socketErr";
	socketErr.appendChild(socketErrMsg);

	document.querySelector("main > .cw-center:first-child").insertAdjacentElement("afterend", socketErr);
}

/** @returns {void} */
function hideError() {
	let socketErr = /** @type {HTMLDivElement} */ (document.getElementById("socketErr"));
	if (socketErr) {
		socketErr.remove();
	}
}
