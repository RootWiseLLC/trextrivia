import { Component } from '@angular/core';
import { RoundState } from '../../model/model';
import { GameStateService } from '../../services/game-state.service';
import { PlayerService } from '../../services/player.service';
import { WebsocketService } from '../../services/websocket.service';

@Component({
    selector: 'app-recv-wager',
    templateUrl: './recv-wager.component.html',
    styleUrls: ['./recv-wager.component.less']
})
export class RecvWagerComponent {
    wagerAmt: string

    constructor(
        private websocket: WebsocketService,
        protected game: GameStateService,
        protected player: PlayerService,
    ) { }

    // The wager range the server enforces. A score below the round's top clue value
    // - including a negative one - still gets the full round maximum to bet with,
    // which is why the range is worth spelling out on screen.
    MinWager(): number {
        return this.game.FinalRound() ? 0 : 5
    }

    MaxWager(): number {
        return Math.max(this.player.Score(), this.roundMax())
    }

    private roundMax(): number {
        switch (this.game.Round()) {
            case RoundState.FirstRound:
                return 1000
            case RoundState.SecondRound:
                return 2000
            default:
                return 0
        }
    }

    handleWager() {
        if (this.wagerAmt != null && this.wagerAmt !== '' && this.player.CanWager()) {
            this.websocket.Send({
                state: this.game.State(),
                wager: this.wagerAmt,
            })
        }
        this.wagerAmt = ''
    }
}
