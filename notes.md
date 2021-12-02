# Table actions

- join: a participant joins the table and chooses his / her player name
- leave: a participant leaves the table
- startRound: opens the table for votes
- vote: a participant submits his / her estimate
- closeRound: when all participants have given their vote or when an admin closes the round manually, no further votes
  are allowed, the cards are revealed
- resetRound: the round is reopened and all votes are reset

# How to add a new action

- [ ] Implement logic in `board`
- [ ] Create Event that triggers the new logic
- [ ] Create Command that checks the new logic and creates the new event(s)
- [ ] Add state related data to `State` and add the mapping logic
- [ ] Add a new `CommandRequest` that is translated to the new Command
- [ ] Add the action name to `mux` and create a variable of the new `CommandRequest`