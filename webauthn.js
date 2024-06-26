class WebAuthnAPI {
	constructor(){}

	async beginRegistration() {
		let path = `/admin/login/webauthn/register:begin`;
		let rsp = await fetch(path, {
			method: 'POST',
		});
		if (!rsp.ok) {
			throw new Error(await rsp.text());
		}
		let options = await rsp.json();
		options.publicKey.challenge = await this.base64decode(options.publicKey.challenge);
		options.publicKey.user.id = await this.base64decode(options.publicKey.user.id);
		return options;
	}
	async finishRegistration(options) {
		options.rawId = await this.base64encode(options.rawId);
		options.response.clientDataJSON = await this.base64encode(options.response.clientDataJSON);
		options.response.attestationObject = await this.base64encode(options.response.attestationObject);

		let path = `/admin/login/webauthn/register:finish`;
		let rsp = await fetch(path, {
			method: 'POST',
			body: JSON.stringify(options),
		});
		if (!rsp.ok) {
			throw new Error(await rsp.text());
		}
	}
	async beginLogin() {
		let path = `/admin/login/webauthn/login:begin`;
		let rsp = await fetch(path, {
			method: 'POST',
		});
		if (!rsp.ok) {
			throw new Error(await rsp.text());
		}
		let options = await rsp.json();
		options.publicKey.challenge = await this.base64decode(options.publicKey.challenge);
		return options;
	}
	// challenge: 只是用来查找 session 的，肯定不会用作服务端真实 challenge。想啥呢？
	async finishLogin(options, challenge) {
		let path = `/admin/login/webauthn/login:finish?challenge=${challenge}`;
		options.rawId = await this.base64encode(options.rawId);
		options.response.authenticatorData = await this.base64encode(options.response.authenticatorData);
		options.response.clientDataJSON = await this.base64encode(options.response.clientDataJSON);
		options.response.signature = await this.base64encode(options.response.signature);
		options.response.userHandle = await this.base64encode(options.response.userHandle);
		let rsp = await fetch(path, {
			method: 'POST',
			body: JSON.stringify(options),
		});
		if (!rsp.ok) {
			throw new Error(await rsp.text());
		}
	}

	async base64encode(arr) {
		let path = '/admin/login/webauthn/base64:encode';
		let body = JSON.stringify(Array.from(new Uint8Array(arr)));
		let rsp = await fetch(path, { method: 'POST', body: body});
		return await rsp.text();
	}
	async base64decode(str) {
		let path = '/admin/login/webauthn/base64:decode';
		let rsp = await fetch(path, { method: 'POST', body: str});
		return new Uint8Array(await rsp.json());
	}
}

class WebAuthn {
	constructor() {
		this.api = new WebAuthnAPI();
	}

	async register() {
		let { publicKey } = await this.api.beginRegistration();
		console.log('开始注册返回：', publicKey);

		let { rp, challenge, user } = publicKey;
		
		let createOptions = {
			challenge: challenge,
			rp: {
				name: rp.name,
				id: rp.id,
			},
			pubKeyCredParams: [
				{
					alg: -8,    // ed25519，我的最爱！❤️
					type: 'public-key',
				},
			],
			user: {
				id: user.id,
				name: user.name,
				displayName: user.displayName,
			},
		};
		console.log('创建参数：', createOptions);

		let credential = null;
		try {
			credential = await navigator.credentials.create({ publicKey: createOptions });
		} catch (e) {
			console.log(e);
			throw e;
		}
		if (!credential) { return null; }
		console.log('本地创建凭证：', credential);

		let finishOptions = {
			id:     credential.id,
			type:   credential.type,
			rawId:  credential.rawId,
			response: {
				clientDataJSON:     credential.response.clientDataJSON,
				attestationObject:  credential.response.attestationObject,
			},
		};
		console.log('结束注册参数：', finishOptions);

		await this.api.finishRegistration(finishOptions);
		console.log('注册成功。');
	}

	async login() {
		let { publicKey } = await this.api.beginLogin();
		console.log('开始登录返回：', publicKey);
		let options = {
			challenge: publicKey.challenge,
		};
		let credential = await navigator.credentials.get({ publicKey: options});
		console.log('本地获取凭证：', credential);

		let finishLoginOptions = {
			id:     credential.id,
			type:   credential.type,
			rawId:  credential.rawId,
			response: {
				clientDataJSON:     credential.response.clientDataJSON,
				authenticatorData:  credential.response.authenticatorData,
				signature:          credential.response.signature,
				userHandle:         credential.response.userHandle,
			},
		};
		console.log('结束登录参数：', finishLoginOptions);
		await this.api.finishLogin(finishLoginOptions, await this.api.base64encode(publicKey.challenge));
		console.log('登录成功。');
	}
}
