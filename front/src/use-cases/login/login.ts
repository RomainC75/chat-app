import { createAction } from "@reduxjs/toolkit";

export const answerSubmitted =
  createAction<AnswerValidation>("ANSWER_SUBMITTED");

export const alreadyValidated = createAction<void>("ALREADY_VALIDATED");

export const login =
  (questionId: string, answerLetter: AnswerLetter): AppThunk<Promise<void>> =>
  async (dispatch, getState, { questionGateway }) => {
    if (getState().answerValidation.data != null) {
      dispatch(alreadyValidated());
      return;
    }
    const answerValidation = await questionGateway.submitAnswer(
      questionId!,
      answerLetter,
    );
    dispatch(answerSubmitted(answerValidation));
  };
