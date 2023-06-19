export const Action = {
	onboard: { action: 'onboard' },
	pause: { action: 'pause' },
	resume: { action: 'resume' },
	set_string: { action: 'set-string' },
	offset_sw: {
		add1: { action: 'offset-sw', value: 1000 },
		sub1: { action: 'offset-sw', value: -1000 },
		add10: { action: 'offset-sw', value: 10000 },
		sub10: { action: 'offset-sw', value: -10000 },
		add60: { action: 'offset-sw', value: 60000 },
		sub60: { action: 'offset-sw', value: -60000 }
	},
	offset_int: {
		home: {
			score_add1: { action: 'offset-int', name: 'home_score', value: 1 },
			score_sub1: { action: 'offset-int', name: 'home_score', value: -1 },
			score_attempt_add1: { action: 'offset-int', name: 'home_scoring_attempts', value: 1 },
			score_attempt_sub1: { action: 'offset-int', name: 'home_scoring_attempts', value: -1 },
			defense_success_add1: { action: 'offset-int', name: 'home_defense_success', value: 1 },
			defense_success_sub1: { action: 'offset-int', name: 'home_defense_success', value: -1 },
			fbg_add1: { action: 'offset-int', name: 'home_fast_break_goals', value: 1 },
			fbg_sub1: { action: 'offset-int', name: 'home_fast_break_goals', value: -1 },
			goal7m_add1: { action: 'offset-int', name: 'home_7m_goals', value: 1 },
			goal7m_sub1: { action: 'offset-int', name: 'home_7m_goals', value: -1 }
		},
		away: {
			score_add1: { action: 'offset-int', name: 'away_score', value: 1 },
			score_sub1: { action: 'offset-int', name: 'away_score', value: -1 },
			score_attempt_add1: { action: 'offset-int', name: 'away_scoring_attempts', value: 1 },
			score_attempt_sub1: { action: 'offset-int', name: 'away_scoring_attempts', value: -1 },
			defense_success_add1: { action: 'offset-int', name: 'away_defense_success', value: 1 },
			defense_success_sub1: { action: 'offset-int', name: 'away_defense_success', value: -1 },
			fbg_add1: { action: 'offset-int', name: 'away_fast_break_goals', value: 1 },
			fbg_sub1: { action: 'offset-int', name: 'away_fast_break_goals', value: -1 },
			goal7m_add1: { action: 'offset-int', name: 'away_7m_goals', value: 1 },
			goal7m_sub1: { action: 'offset-int', name: 'away_7m_goals', value: -1 }
		}
	},
	playout: {
		pregame: {
			countdown: { action: 'playout', name: 'countdown' },
			announcement: { action: 'playout', name: 'announcement' },
			lineup: {
				home: { action: 'playout', name: 'lineup', value: 'home' },
				away: { action: 'playout', name: 'lineup', value: 'away' }
			},
			referees: { action: 'playout', name: 'referees' }
		},
		ingame: {
			permanent_clock: { action: 'playout', name: 'scoreboard' },
			card: {
				home: {
					yellow: { action: 'playout', name: 'home_card', value: 'yellow' },
					white: { action: 'playout', name: 'home_card', value: 'white' },
					blue: { action: 'playout', name: 'home_card', value: 'blue' },
					red: { action: 'playout', name: 'home_card', value: 'red' }
				},
				away: {
					yellow: { action: 'playout', name: 'away_card', value: 'yellow' },
					white: { action: 'playout', name: 'away_card', value: 'white' },
					blue: { action: 'playout', name: 'away_card', value: 'blue' },
					red: { action: 'playout', name: 'away_card', value: 'red' }
				}
			}
		},
		postgame: {
			ff_result: { action: 'playout', name: 'ff_result' },
			statistics: { action: 'playout', name: 'statistics' }
		}
	}
};
