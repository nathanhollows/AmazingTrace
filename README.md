# Amazing Trace

**A QR code based scavenger hint written in Go**

<!-- ABOUT THE PROJECT -->
## About The Project

The Amazing Trace was developed for the [Locals Collegiate Community](https://www.otago.ac.nz/locals/index.html) as an O-week event. The goal of the game was to introduce our first year students to locations around the university. Each team was given a set of three clues for three random locations. They had to solve the clue, get to the location, and scan the poster. Some locations had challenges lead by team members. Upon completing a challenge the team would be allowed to scan the poster, and in some cases be given the opportunity to help themselves (solve the next three clues automatically) or hinder others (shuffle their clues randomly).

The project was a huge hit with great turn out. I'd really like to run this again and iron out some of the wrinkles.

## Game play

Users are given a team code, which are randomly generated when the game is launched. The team code is written on a card and given out. The card outlines very briefly how to play.

<img src="https://user-images.githubusercontent.com/13064427/110064102-599c8980-7dd1-11eb-9f4f-e29f2d64e906.png" width="50%">
<img src="https://user-images.githubusercontent.com/13064427/110064104-5a352000-7dd1-11eb-8404-98fc87c9b88d.png" width="50%">

The homepage asks for the team code which is then remembered. The code is required for all future clues.

<img src="https://user-images.githubusercontent.com/13064427/110063888-dd09ab00-7dd0-11eb-9875-33344a40ed5d.png" width="50%">

These posters were hidden around the university. For people that didn't have a QR scanner, the corresponding URL was printed on every poster.

<img src="https://user-images.githubusercontent.com/13064427/110063919-ee52b780-7dd0-11eb-8f6b-047893e76048.png" width="50%">

The game used a timer to mark the end of the game. When players scanned a clue the amount of time left was at the top of the page. The game instructed all teams to return when time was up, and the admin panel showed all team scores, and who won.

## What I would have changed

**1. More time to play**

All teams wanted more time to play. There were 25 clues but the winning team only found 16 of them. We were running with time constraints and there wasn't much we could do about it.

**2. A full team play through ahead of time**

The first time we ran the game was with our first year students. We essentially only tested the project in production (apart from my own testing). It would have got the leaders team on the same page and made the game run much smoother.

~~**3. Dyanmically add teams**~~

~~The game currently generates 50 teams codes and only adds them to the admin panel as the teams log in. I was struggling with the Go language and this was the easiest way around it given the timeframe I had to complete it. As long as there are fewer than 50 teams it would work well.~~

Done!

**4. Dynamically adjust end time**

The end time was hardcoded into the game. It didn't affect the game, it just felt like a dirty approach.

<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.


<!-- CONTACT -->
## Contact

Nathan Hollows - me@nathanhollows.com
