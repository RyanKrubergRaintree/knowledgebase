export {};

interface TestCaseAssert {
	/** Throws error with msg if not okay */
	true: (ok: boolean, msg: string) => void;
	/** Throws error with message */
	fail: (err: string) => void;
	equal: (actual: any, expect: any, msg: string) => void;
}

interface getClassListResult {
	/** add token to className */
	add: (token) => void;
	/** remove token to className */
	remove: (token) => void;
	/** toggle token to className */
	toggle: (token) => void;
	/** get if token is in className */
	contains: (token) => void;
}

/**
 * Reloads files on change
 * https://github.com/raintreeinc/livepkg/blob/master/reloader.js.go
 */
interface Reloader {
	/** time in **milliseconds** @default 2000 */
	ReloadAfter: number;
	loading: object;
	unloaded: [];
	onchange?: (change?: any) => void;
	Change?: (change?: any) => void;
}

declare global {
	// --- Types for [livepkg](https://github.com/raintreeinc/livepkg/tree/master)
	/**
	 * Used for defining javascript packages -- the raintree way
	 *
	 * {@link https://github.com/raintreeinc/livepkg/blob/master/package.js.go}
	 */
	function package(name: string, setup: (namespace: object) => object | undefined | void): void;

	/**
	 * Used for defining JavaScript (js) dependencies that the current javascript file requires
	 * Can be later used by the backend to generate js requirement
	 * @example `<script src="${file_path}" type="text/javascript" >`
	 *
	 * {@link https://github.com/raintreeinc/livepkg/blob/master/file.go}
	 *
	 * @param file_path path to the file
	 */
	function depends(file_path: string): void;

	/**
	 * Reloads files on change
	 *
	 * {@link https://github.com/raintreeinc/livepkg/blob/master/reloader.js.go}
	 */
	const Reloader: Reloader;

	// --- Types for globals declared in client\assets\js\global.js ---
	/**
	 * Gets element data attribute.
	 *
	 * Tries `el.dataset[name]` then `el.getAttribute("data-" + name)`
	 */
	function GetDataAttribute(el: HTMLElement, name: string): string;
	/**
	 * Generates random string of (16-3)*2=26 characters
	 *
	 * Can be 25 at times
	 *
	 * @example
	 * GenerateID() -> "2a8fe8217e6cd4b9e5016cbb6b"
	 * GenerateID() -> "9db30e6c439caf74717849ffb"
	 */
	function GenerateID(): string;

	/** Runs a test case */
	function TestCase(casename: string, runcase: (assert: TestCaseAssert) => void): void;

	function getClassList(el: HTMLElement): getClassListResult;
}
