import { Component, OnInit } from '@angular/core'
import { Router } from '@angular/router'
import { Observer } from 'rxjs'
import { ServerUnavailableMsg, User } from '../model/model'
import { ApiService } from '../services/api.service'
import { AuthService } from '../services/auth.service'
import { JwtService } from '../services/jwt.service'
import { ModalService } from '../services/modal.service'
import { environment } from '../../environments/environment'

@Component({
	selector: 'app-join',
	templateUrl: './join.component.html',
	styleUrls: ['./join.component.less'],
})
export class JoinComponent implements OnInit {
	protected user: User
	protected joinCode: string = ''
	protected nameInvalid: boolean = false
	protected showDedication: boolean = environment.showBirthdayMessage
	protected oneRoundChecked: boolean = false
	protected twoRoundChecked: boolean = true
	protected penaltyChecked: boolean = true

	constructor(
		private router: Router,
		private jwt: JwtService,
		private api: ApiService,
		protected modal: ModalService,
		private auth: AuthService,
	) { }

	ngOnInit() {
		this.auth.user.subscribe(user => {
			this.user = user
		})
	}

	private joinResp(): Partial<Observer<any>> {
		return {
			next: (resp: any) => {
				this.auth.RememberName(this.user.name)
				this.jwt.SetJWT(resp.token)
				this.router.navigate([`/game/${resp.game.name}`])
			},
			error: (err: any) => {
				let msg = err.status != 0 ? err.error.message : ServerUnavailableMsg;
				this.modal.displayMessage(msg)
			}
		}
	}

	private showInvalidName() {
		document.getElementById('player-name-join')?.focus()
		this.nameInvalid = true
		setTimeout(() => {
			this.nameInvalid = false
		}, 1000)
	}

	setRounds(rounds: number) {
		this.oneRoundChecked = rounds == 1
		this.twoRoundChecked = rounds == 2
	}

	createPrivateGame(bots: number) {
		if (this.user.name == '') {
			this.showInvalidName()
			return
		}
		this.api.CreatePrivateGame(
			this.user.name, this.user.imgUrl, this.user.email,
			bots, this.twoRoundChecked, this.penaltyChecked,
			30, 30, 15, 30, [], []
		).subscribe(this.joinResp())
	}

	joinGameByCode() {
		if (this.user.name == '') {
			this.showInvalidName()
			return
		}
		this.api.JoinGameByCode(this.user.name, this.user.imgUrl, this.user.email, this.joinCode).subscribe(this.joinResp())
	}

	joinPublicGame() {
		if (this.user.name == '') {
			this.showInvalidName()
			return
		}
		this.api.JoinPublicGame(this.user.name, this.user.imgUrl, this.user.email).subscribe(this.joinResp())
	}

	rejoin() {
		this.api.GetPlayerGame().subscribe({
			next: (resp: any) => {
				this.router.navigate([`/game/${resp.game.name}`])
			},
			error: (err: any) => {
				let msg = err.status != 0 ? err.error.message : ServerUnavailableMsg
				this.modal.displayMessage(msg)
			},
		})
	}
}
