#include <iostream>
#include <fstream>
#include <sstream>
#include <llvm/IR/LLVMContext.h>
#include <llvm/IR/IRBuilder.h>
#include <llvm/IR/Module.h>
#include <llvm/IR/Value.h>
#include <llvm/Support/raw_ostream.h>
#include <llvm/IR/Verifier.h>
#include <llvm/Passes/PassBuilder.h>
#include <llvm/Passes/OptimizationLevel.h>
#include <map>
#include <string>
#include <memory>
#include <stack>

using namespace llvm;

static std::unique_ptr<LLVMContext> TheContext;
static std::unique_ptr<IRBuilder<>> Builder;
static std::unique_ptr<Module> TheModule;

Value *getTapeIndex(GlobalVariable *tapeIndex)
{
    return Builder->CreateLoad(Type::getInt32Ty(*TheContext), tapeIndex, "currentIndex");
}

Value *getTapePtr(ArrayType *TapeType, GlobalVariable *Tape, Value *currentIndex)
{
    return Builder->CreateGEP(TapeType, Tape, {ConstantInt::get(Type::getInt8Ty(*TheContext), 0), currentIndex}, "tapePtr");
}

void compile(std::string program)
{
    TheContext = std::make_unique<LLVMContext>();
    TheModule = std::make_unique<Module>("my cool bfcompiler", *TheContext);
    Builder = std::make_unique<IRBuilder<>>(*TheContext);

    ArrayType *TapeType = ArrayType::get(Type::getInt8Ty(*TheContext), 1024 * 1024 * 64);
    GlobalVariable *Tape = new GlobalVariable(*TheModule, TapeType, false, GlobalValue::PrivateLinkage,
                                              Constant::getNullValue(TapeType), "tape");

    GlobalVariable *tapeIndex = new GlobalVariable(*TheModule, Type::getInt32Ty(*TheContext), false, GlobalValue::PrivateLinkage,
                                                   ConstantInt::get(Type::getInt32Ty(*TheContext), (1024 * 1024 * 64) / 2), "tapePointer");

    std::stack<std::pair<BasicBlock *, BasicBlock *>> loopStack;

    FunctionType *putcharType = FunctionType::get(Type::getVoidTy(*TheContext), {Type::getInt8Ty(*TheContext)}, false);
    Function *putcharFunc = Function::Create(putcharType, Function::ExternalLinkage, "my_putchar", TheModule.get());

    FunctionType *getcharType = FunctionType::get(Type::getInt8Ty(*TheContext), {}, false);
    Function *getcharFunc = Function::Create(getcharType, Function::ExternalLinkage, "my_getchar", TheModule.get());

    FunctionType *mainFuncType = FunctionType::get(Type::getInt32Ty(*TheContext), false);
    Function *mainFunction = Function::Create(mainFuncType, Function::ExternalLinkage, "main", TheModule.get());
    BasicBlock *entryBlock = BasicBlock::Create(*TheContext, "entry", mainFunction);
    Builder->SetInsertPoint(entryBlock);

    for (char ch : program)
    {
        switch (ch)
        {
        case '>':
        {
            Value *currentIndex = getTapeIndex(tapeIndex);
            Value *newVal = Builder->CreateAdd(currentIndex, ConstantInt::get(Type::getInt32Ty(*TheContext), 1), "newVal");
            Builder->CreateStore(newVal, tapeIndex);
            break;
        }
        case '<':
        {
            Value *currentIndex = getTapeIndex(tapeIndex);
            Value *newVal = Builder->CreateSub(currentIndex, ConstantInt::get(Type::getInt32Ty(*TheContext), 1), "newVal");
            Builder->CreateStore(newVal, tapeIndex);
            break;
        }
        case '+':
        {
            Value *currentIndex = getTapeIndex(tapeIndex);
            Value *ptr = getTapePtr(TapeType, Tape, currentIndex);
            Value *currentVal = Builder->CreateLoad(Type::getInt8Ty(*TheContext), ptr, "currentVal");
            Value *newVal = Builder->CreateAdd(currentVal, ConstantInt::get(Type::getInt8Ty(*TheContext), 1), "newVal");
            Builder->CreateStore(newVal, ptr);
            break;
        }
        case '-':
        {
            Value *currentIndex = getTapeIndex(tapeIndex);
            Value *ptr = getTapePtr(TapeType, Tape, currentIndex);
            Value *currentVal = Builder->CreateLoad(Type::getInt8Ty(*TheContext), ptr, "currentVal");
            Value *newVal = Builder->CreateSub(currentVal, ConstantInt::get(Type::getInt8Ty(*TheContext), 1), "newVal");
            Builder->CreateStore(newVal, ptr);
            break;
        }
        break;
        case '.':
        {
            Value *currentIndex = getTapeIndex(tapeIndex);
            Value *ptr = getTapePtr(TapeType, Tape, currentIndex);
            Value *currentVal = Builder->CreateLoad(Type::getInt8Ty(*TheContext), ptr, "currentVal");
            Builder->CreateCall(putcharFunc, currentVal);
            break;
        }
        case ',':
        {
            Value *inputVal = Builder->CreateCall(getcharFunc);
            Value *currentIndex = Builder->CreateLoad(Type::getInt8Ty(*TheContext), tapeIndex, "currentIndex");
            Value *ptr = getTapePtr(TapeType, Tape, currentIndex);
            Builder->CreateStore(inputVal, ptr);
            break;
        }
        case '[':
        {
            BasicBlock *loopCond = BasicBlock::Create(*TheContext, "loopCond", mainFunction);
            BasicBlock *loopBody = BasicBlock::Create(*TheContext, "loopBody", mainFunction);
            BasicBlock *afterLoop = BasicBlock::Create(*TheContext, "afterLoop", mainFunction);
            Builder->CreateBr(loopCond);
            Builder->SetInsertPoint(loopCond);

            Value *currentIndex = getTapeIndex(tapeIndex);
            Value *ptr = getTapePtr(TapeType, Tape, currentIndex);
            Value *cellValue = Builder->CreateLoad(Type::getInt8Ty(*TheContext), ptr, "cellValue");
            Value *isZero = Builder->CreateICmpEQ(cellValue, ConstantInt::get(Type::getInt8Ty(*TheContext), 0), "isZero");
            Builder->CreateCondBr(isZero, afterLoop, loopBody);
            Builder->SetInsertPoint(loopBody);
            loopStack.push({loopCond, afterLoop});
            break;
        }
        case ']':
        {
            auto pair = loopStack.top();
            auto loopCond = pair.first;
            auto afterLoop = pair.second;
            loopStack.pop();
            Builder->CreateBr(loopCond);
            Builder->SetInsertPoint(afterLoop);
            break;
        }

        default:
            break; // do nothing
        }
    }

    Builder->CreateRet(ConstantInt::get(Type::getInt32Ty(*TheContext), 0));
}

int main(int argc, char *argv[])
{
    if (argc <= 1)
    {
        std::cout << "Usage: ./bf <bf-program.b>";
        return 0;
    }

    std::ifstream file(argv[1]);
    if (!file)
    { // Check if the file was successfully opened
        std::cerr << "Could not open the file" << argv[1] << std::endl;
        return 1;
    }

    std::stringstream buffer;
    buffer << file.rdbuf();

    std::string program = buffer.str();

    compile(program);

    auto res = verifyModule(*TheModule, &llvm::errs());
    assert(!res); // Check that there are no verification errors

    // taken from LLVM docs
    // https://rocm.docs.amd.com/projects/llvm-project/en/latest/LLVM/llvm/html/NewPassManager.html
    // Create the analysis managers.
    // These must be declared in this order so that they are destroyed in the
    // correct order due to inter-analysis-manager references.
    LoopAnalysisManager LAM;
    FunctionAnalysisManager FAM;
    CGSCCAnalysisManager CGAM;
    ModuleAnalysisManager MAM;

    // Create the new pass manager builder.
    // Take a look at the PassBuilder constructor parameters for more
    // customization, e.g. specifying a TargetMachine or various debugging
    // options.
    PassBuilder PB;

    // Register all the basic analyses with the managers.
    PB.registerModuleAnalyses(MAM);
    PB.registerCGSCCAnalyses(CGAM);
    PB.registerFunctionAnalyses(FAM);
    PB.registerLoopAnalyses(LAM);
    PB.crossRegisterProxies(LAM, FAM, CGAM, MAM);

    // Create the pass manager.
    // This one corresponds to a typical -O2 optimization pipeline.
    ModulePassManager MPM = PB.buildPerModuleDefaultPipeline(OptimizationLevel::O3);

    // Optimize the IR!
    MPM.run(*TheModule, MAM);

    TheModule->print(llvm::outs(), nullptr);

    return 0;
}