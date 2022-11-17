package votesy.results.Controller;

import java.util.ArrayList;

import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;

import votesy.results.AzureTableService;
import votesy.results.Question;
import votesy.results.Vote;

@Controller
public class WebController {
    @RequestMapping(value = "/")
    public String index(@RequestParam(name="id", required=false, defaultValue="0") String name, Model model) {
        AzureTableService ats = new AzureTableService();
        Question currentQuestion = ats.getCurrentQuestion();
        ArrayList<Vote> votes = ats.getVotesForQuestion(currentQuestion);
        String voteUrl = System.getenv("voteUrl");

        System.out.println(voteUrl);

        model.addAttribute("question", currentQuestion.text);
        model.addAttribute("answer1", currentQuestion.answer1Text);
        model.addAttribute("answer2", currentQuestion.answer2Text);

        model.addAttribute("votes1", votes.get(0).VoteCount);
        model.addAttribute("votes2", votes.get(1).VoteCount);

        model.addAttribute("voteUrl", voteUrl);
        return "index";
    }
}
