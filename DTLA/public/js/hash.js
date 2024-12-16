"use strict";

/**
 * All the elements where user input might cause an error
 * @typedef {(HTMLElement | HTMLFieldSetElement | HTMLTextAreaElement | HTMLInputElement)} HashFormElem
 */

/** @type {HTMLFormElement} */
let hashForm;

/** @type {HTMLFieldSetElement} */
let hashFormFieldset;

/** @type {NodeListOf<HTMLInputElement>} */
let hashFormFuncRadios;

/** @type {HTMLTextAreaElement} */
let hashFormIn;

/** @type {HTMLInputElement} */
let hashFormOut;

/** @type {HTMLButtonElement} */
let hashFormSubmit;

window.addEventListener("load", () => {
	hashForm =
		/** @type {HTMLFormElement} */
		(document.getElementById("hash-form"));
	if (!hashForm)
		throw Error();

	hashFormFieldset =
		/** @type {HTMLFieldSetElement} */
		(document.getElementById("hash-form-funcs"));

	hashFormFuncRadios =
		/** @type {NodeListOf<HTMLInputElement>} */
		(document.getElementsByName("hash-form-func"));
	if (hashFormFuncRadios.length == 0)
		throw Error();

	hashFormIn =
		/** @type {HTMLTextAreaElement} */
		(document.getElementById("hash-form-input"));
	if (!hashFormIn)
		throw Error();

	hashFormOut =
		/** @type {HTMLInputElement} */
		(document.getElementById("hash-form-output"));
	if (!hashFormOut)
		throw Error();

	hashFormSubmit =
		/** @type {HTMLButtonElement} */
		(document.getElementById("hash-form-submit"));
	if (!hashFormSubmit)
		throw Error();

	hashFormSubmit.addEventListener("click", hashHandleSubmit);
});

function hashHandleSubmit() {
	let hashFunc = /** @type {HTMLInputElement} */(document.querySelector("input[name='hash-form-func']:checked"));
	if (!hashFunc) {
		hashShowError(hashFormFieldset, "Nav izvēlēta hash funkcija.");
		return;
	}
	hashHideError();

	if (hashFormIn.value == "") {
		hashShowError(hashFormIn, "Nav ievadīts ievades teksts");
		return;
	}
	hashHideError
	
	getHash(hashFunc.value, hashFormIn.value).then(
		(hashHex) => {
			hashFormOut.value = hashHex;
		},
		() => {
			hashFormOut.value = "Notika kļūme";
		});
}

/**
 * Return a hash of the plain text UTF8 string using the passed in algorithm
 * @param {AlgorithmIdentifier} algorithm
 * @param {string} plaintext
 * @returns {Promise<string>}
 */
async function getHash(algorithm, plaintext) {
	return Array.from(new Uint8Array(await crypto.subtle.digest(algorithm, new TextEncoder().encode(plaintext))))
		.map((b) => b.toString(16).padStart(2, "0")).join("");
}

/**
 * Highlight passed in elem and err msg elem in red, choose style.
 * Insert a new error msg elem if not exists in grid row 3, col [1;3], max-content or width: 100%.
 * @param {HashFormElem} elemInvalid
 * @param {string} msg 
 */
function hashShowError(elemInvalid, msg) {
	elemInvalid.style.outline = "none";
	elemInvalid.style.border = "1px solid red";

	let elemErrMsg = /** @type {HTMLParagraphElement} */ (document.getElementById("hash-form-error"));
	if (elemErrMsg) {
		elemErrMsg.innerText = msg;
		return;
	}

	elemErrMsg = document.createElement("p");
	elemErrMsg.id = "hash-form-error";
	elemErrMsg.innerText = msg;
	hashFormSubmit.insertAdjacentElement("afterend", elemErrMsg);
}

/** Hides any error messages and other styling that's relevant if there was an error for hashForm. */
function hashHideError() {
	hashForm.childNodes.forEach((/** @type {HTMLElement} */elem) => {
		if (elem.style != undefined) {
			elem.style.outline = null;
			elem.style.border = null;
		}
	})

	let elemErrMsg = /** @type {HTMLParagraphElement} */ (document.getElementById("hash-form-error"));
	if (elemErrMsg) {
		elemErrMsg.remove();
	}
}
